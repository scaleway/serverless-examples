package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/scaleway/scaleway-sdk-go/api/block/v1"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	envOrgID        = "SCW_DEFAULT_ORGANIZATION_ID"
	envAccessKey    = "SCW_ACCESS_KEY"
	envSecretKey    = "SCW_SECRET_KEY"
	envProjectID    = "SCW_DEFAULT_PROJECT_ID"
	envInstanceZone = "SCW_ZONE"
	envInstanceID   = "INSTANCE_ID"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	if err := checkEnvVars(); err != nil {
		return err
	}

	slog.Info("creating snapshot of instance...")

	// Create a Scaleway client with credentials from environment variables.
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv(envOrgID)),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv(envAccessKey), os.Getenv(envSecretKey)),

		// Get more about our availability
		// zones at https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
		scw.WithDefaultRegion(scw.RegionFrPar),

		scw.WithDefaultZone(scw.Zone(os.Getenv(envInstanceZone))),
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Create SDK objects for Scaleway Instance product
	instanceAPI := instance.NewAPI(client)
	blockAPI := block.NewAPI(client)

	ctx := context.Background()

	if err := createSnapshots(ctx, instanceAPI, blockAPI); err != nil {
		return err
	}

	return nil
}

func createSnapshots(ctx context.Context, instanceAPI *instance.API, blockAPI *block.API) error {
	zone := scw.Zone(os.Getenv(envInstanceZone))
	instanceID := os.Getenv(envInstanceID)

	gotInstance, err := instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: instanceID,
		Zone:     zone,
	}, scw.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error while getting instance: %w", err)
	}

	for _, volume := range gotInstance.Server.Volumes {
		slog.Info("getting volume informations", "volume_id", volume.ID)

		volumeHydrated, err := blockAPI.GetVolume(&block.GetVolumeRequest{
			VolumeID: volume.ID,
			Zone:     zone,
		})
		if err != nil {
			return fmt.Errorf("erro while reading volume informations: %w", err)
		}

		slog.Info("creating snapshot for volume", "volume_id", volume.ID)
		snapshotResp, err := blockAPI.CreateSnapshot(&block.CreateSnapshotRequest{
			Zone:      zone,
			VolumeID:  volume.ID,
			ProjectID: os.Getenv(envProjectID),
			Name:      fmt.Sprintf("snapshot-%s-%s", volumeHydrated.Name, time.Now().Format("2006-01-02")),
		}, scw.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("error while creating snapshot: %w", err)
		}

		slog.Info("created snapshot", "snapshot_id", snapshotResp.ID)
	}

	return nil
}

func checkEnvVars() error {
	mandatoryVariables := []string{envOrgID, envAccessKey, envSecretKey, envInstanceID, envInstanceZone, envProjectID}

	for _, v := range mandatoryVariables {
		if os.Getenv(v) == "" {
			return fmt.Errorf("missing environment variable: %s", v)
		}
	}
	return nil
}
