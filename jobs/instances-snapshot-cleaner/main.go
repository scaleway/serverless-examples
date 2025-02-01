package main

import (
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
	envProjectID = "SCW_PROJECT_ID"
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
		panic(err)
	}
}

func cleanSnapshots(days int, instanceAPI *instance.API) error {
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

	for _, snapshot := range snapshotsList.Snapshots {
		if snapshot.State == instance.SnapshotStateAvailable && (currentTime.Sub(*snapshot.CreationDate).Hours()/hoursPerDay) > float64(days) {
			fmt.Printf("\nDeleting snapshot <%s>:%s created at: %s\n", snapshot.ID, snapshot.Name, snapshot.CreationDate.Format(time.RFC3339))

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

func init() {
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envZone, envProjectID}

	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}
