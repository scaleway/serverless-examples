package main

import (
	"fmt"
	"os"
	"time"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	envOrgID        = "SCW_DEFAULT_ORGANIZATION_ID"
	envAccessKey    = "SCW_ACCESS_KEY"
	envSecretKey    = "SCW_SECRET_KEY"
	envInstanceZone = "SCW_ZONE"
	envInstanceID   = "INSTANCE_ID"
)

func main() {
	fmt.Println("creating snapshot of instance...")

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

	if err := createSnapshots(instanceAPI); err != nil {
		panic(err)
	}
}

func createSnapshots(instanceAPI *instance.API) error {
	gotInstance, err := instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: os.Getenv(envInstanceID),
		Zone:     scw.Zone(os.Getenv(envInstanceZone)),
	})
	if err != nil {
		return fmt.Errorf("error while getting instance %w", err)
	}

	now := time.Now().Format(time.DateOnly)

	for _, volume := range gotInstance.Server.Volumes {
		snapshotName := fmt.Sprintf("snap-vol-%s-%s-%s",
			volume.VolumeType.String(),
			now,
			os.Getenv(envInstanceZone))

		snapshotResp, err := instanceAPI.CreateSnapshot(&instance.CreateSnapshotRequest{
			Name:       snapshotName,
			VolumeID:   &volume.ID,
			VolumeType: instance.SnapshotVolumeType(volume.VolumeType),
			Zone:       scw.Zone(os.Getenv(envInstanceZone)),
		})
		if err != nil {
			return fmt.Errorf("error while creating snapshot %w", err)
		}

		fmt.Println("created snapshot ", snapshotResp.Snapshot.ID)
	}

	return nil
}

func init() {
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envInstanceID, envInstanceZone}

	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}
