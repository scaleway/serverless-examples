package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// envDeleteAfter name of env variable to deleter older images.
	envDeleteAfter = "SCW_DELETE_AFTER_DAYS"

	// defaultDaysDeleteAfter is the default days value for older images to be deleted.
	defaultDaysDeleteAfter = int(90)
)

func main() {
	fmt.Println("cleaning instances snapshots...")

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
		panic(err)
	}

	// Create SDK objects for Scaleway Instance product
	instanceAPI := instance.NewAPI(client)

	deleteAfterDays := defaultDaysDeleteAfter

	deleteAfterDaysVar := os.Getenv(envDeleteAfter)

	if deleteAfterDaysVar != "" {
		deleteAfterDays, err = strconv.Atoi(deleteAfterDaysVar)
		if err != nil {
			panic(err)
		}
	}

	if err := cleanSnapshots(deleteAfterDays, instanceAPI); err != nil {
		var precondErr *scw.PreconditionFailedError

		if errors.As(err, &precondErr) {
			fmt.Println("\nExtracted Error Details:")
			fmt.Println("Precondition:", precondErr.Precondition)
			fmt.Println("Help Message:", precondErr.HelpMessage)

			// Decode RawBody (if available)
			if len(precondErr.RawBody) > 0 {
				var parsedBody map[string]interface{}
				if json.Unmarshal(precondErr.RawBody, &parsedBody) == nil {
					fmt.Println("RawBody (Decoded):", parsedBody)
				} else {
					fmt.Println("RawBody (Raw):", string(precondErr.RawBody))
				}
			}
		}
		panic(err)
	}
}

// cleanSnapshots when called will clean snapshots in the project (if specified)
// that are older than the number of days.
func cleanSnapshots(days int, instanceAPI *instance.API) error {
	// Get the list of all snapshots
	// TODO: use block api here?
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
			fmt.Printf("\nDeleting snapshot <%s>:%s created at: %s\n", snapshot.ID, snapshot.Name, snapshot.CreationDate.Format(time.RFC3339))

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
