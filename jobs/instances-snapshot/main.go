package main

import (
	"fmt"
	"os"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func main() {
	fmt.Println("creating snapshot of instance...")

	// Create a Scaleway client with credentials from environment variables.
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv("SCW_DEFAULT_ORGANIZATION_ID")),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv("SCW_ACCESS_KEY"), os.Getenv("SCW_SECRET_KEY")),

		// Get more about our availability zones at https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
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
		ServerID: os.Getenv("INSTANCE_ID"),
		Zone:     scw.Zone(os.Getenv("INSTANCE_ZONE")),
	})
	if err != nil {
		return fmt.Errorf("error while getting instance %w", err)
	}

	for _, volume := range gotInstance.Server.Volumes {
		snapshotResp, err := instanceAPI.CreateSnapshot(&instance.CreateSnapshotRequest{
			Name:       volume.Name + RandomString(4),
			VolumeID:   &volume.ID,
			VolumeType: instance.SnapshotVolumeTypeBSSD,
			Zone:       scw.Zone(os.Getenv("INSTANCE_ZONE")),
		})
		if err != nil {
			return fmt.Errorf("error while creating snapshopt %w", err)
		}
		fmt.Println("created snapshot ", snapshotResp.Snapshot.ID)
	}

	return nil
}

func init() {
	if os.Getenv("SCW_DEFAULT_ORGANIZATION_ID") == "" {
		panic("missing SCW_DEFAULT_ORGANIZATION_ID")
	}

	if os.Getenv("SCW_ACCESS_KEY") == "" {
		panic("missing SCW_ACCESS_KEY")
	}

	if os.Getenv("SCW_SECRET_KEY") == "" {
		panic("missing SCW_SECRET_KEY")
	}

	if os.Getenv("INSTANCE_ID") == "" {
		panic("missing INSTANCE_ID")
	}

	if os.Getenv("INSTANCE_ZONE") == "" {
		panic("missing INSTANCE_ZONE")
	}
}
