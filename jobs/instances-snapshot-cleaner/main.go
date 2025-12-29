package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	envOrgID     = "SCW_DEFAULT_ORGANIZATION_ID"
	envAccessKey = "SCW_ACCESS_KEY"
	envSecretKey = "SCW_SECRET_KEY"
	envProjectID = "SCW_DEFAULT_PROJECT_ID"
	envZone      = "SCW_ZONE"
	envLogLevel  = "LOG_LEVEL"

	// envDeleteAfter name of env variable to deleter older images.
	envDeleteAfter = "SCW_DELETE_AFTER_DAYS"

	// defaultDaysDeleteAfter is the default days value for older images to be deleted.
	defaultDaysDeleteAfter = int(90)
)

var logger *slog.Logger

func main() {
	// Initialize structured logger
	logLevel := slog.LevelInfo
	if lvl := os.Getenv(envLogLevel); lvl != "" {
		switch lvl {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		}
	}

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	logger = slog.New(h)
	slog.SetDefault(logger)

	ctx := context.Background()
	logger.InfoContext(ctx, "starting instances snapshots cleaner")

	// Create a Scaleway client with credentials from environment variables.
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv(envOrgID)),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv(envAccessKey), os.Getenv(envSecretKey)),

		// Get more about our availability
		// zones at https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	if err != nil {
		logger.ErrorContext(ctx, "failed to create scaleway client", "error", err)
		panic(err)
	}

	// Create SDK objects for Scaleway Instance product
	instanceAPI := instance.NewAPI(client)

	deleteAfterDays := defaultDaysDeleteAfter

	deleteAfterDaysVar := os.Getenv(envDeleteAfter)

	if deleteAfterDaysVar != "" {
		deleteAfterDays, err = strconv.Atoi(deleteAfterDaysVar)
		if err != nil {
			logger.ErrorContext(ctx, "failed to parse delete after days", "value", deleteAfterDaysVar, "error", err)
			panic(err)
		}
	}

	logger.InfoContext(ctx, "cleaning snapshots", "delete_after_days", deleteAfterDays)

	if err := cleanSnapshotsWithLogging(ctx, deleteAfterDays, instanceAPI); err != nil {
		var precondErr *scw.PreconditionFailedError

		if errors.As(err, &precondErr) {
			logger.ErrorContext(ctx, "scaleway precondition failed",
				"precondition", precondErr.Precondition,
				"help_message", precondErr.HelpMessage)

			// Decode RawBody (if available)
			if len(precondErr.RawBody) > 0 {
				var parsedBody map[string]interface{}
				if json.Unmarshal(precondErr.RawBody, &parsedBody) == nil {
					logger.ErrorContext(ctx, "scaleway error raw body", "body", parsedBody)
				} else {
					logger.ErrorContext(ctx, "scaleway error raw body", "body", string(precondErr.RawBody))
				}
			}
		}
		logger.ErrorContext(ctx, "failed to clean snapshots", "error", err)
		panic(err)
	}

	logger.InfoContext(ctx, "successfully cleaned snapshots")
}

// cleanSnapshotsWithLogging when called will clean snapshots in the project (if specified)
// that are older than the number of days.
func cleanSnapshotsWithLogging(ctx context.Context, days int, instanceAPI *instance.API) error {
	// Get the list of all snapshots
	snapshotsList, err := instanceAPI.ListSnapshots(&instance.ListSnapshotsRequest{
		Zone:    scw.Zone(os.Getenv(envZone)),
		Project: scw.StringPtr(os.Getenv(envProjectID)),
	},
		scw.WithAllPages())
	if err != nil {
		return fmt.Errorf("error while listing snapshots %w", err)
	}

	const hoursPerDay = 24

	currentTime := time.Now()

	// For each snapshot, check conditions
	for _, snapshot := range snapshotsList.Snapshots {
		// Check if snapshot is in ready state and if it's older than the number of days definied.
		if snapshot.State == instance.SnapshotStateAvailable && (currentTime.Sub(*snapshot.CreationDate).Hours()/hoursPerDay) > float64(days) {
			logger.InfoContext(ctx, "deleting snapshot",
				"id", snapshot.ID,
				"name", snapshot.Name,
				"created_at", snapshot.CreationDate.Format(time.RFC3339))

			// Delete snapshot found.
			err := instanceAPI.DeleteSnapshot(&instance.DeleteSnapshotRequest{
				SnapshotID: snapshot.ID,
				Zone:       snapshot.Zone,
			})
			if err != nil {
				return fmt.Errorf("error while deleting snapshot: %w", err)
			}
		}
	}

	return nil
}

// Check for mandatory variables before starting to work.
func init() {
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envZone, envProjectID}

	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}
