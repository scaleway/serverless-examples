package main

import (
	"fmt"
	"os"

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

	blockAPI := block.NewAPI(client)

	if err := createSnapshots(instanceAPI, blockAPI); err != nil {
		panic(err)
	}
}

func createSnapshots(instanceAPI *instance.API, blockAPI *block.API) error {
	gotInstance, err := instanceAPI.GetServer(&instance.GetServerRequest{
		ServerID: os.Getenv(envInstanceID),
		Zone:     scw.Zone(os.Getenv(envInstanceZone)),
	})
	if err != nil {
		return fmt.Errorf("error while getting instance %w", err)
	}

	for _, volume := range gotInstance.Server.Volumes {
		snapshotResp, err := blockAPI.CreateSnapshot(&block.CreateSnapshotRequest{
			Zone:      scw.Zone(os.Getenv(envInstanceZone)),
			VolumeID:  volume.ID,
			ProjectID: os.Getenv(envProjectID),
		})
		if err != nil {
			return fmt.Errorf("error while creating snapshot %w", err)
		}

		fmt.Println("created snapshot ", snapshotResp.ID)
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
