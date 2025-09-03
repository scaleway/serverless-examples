package main

import (
	"fmt"
	"os"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Environment variable constants used to configure the Scaleway API client.
// These must be set in the environment for the application to authenticate and interact with Scaleway services.
const (
	// envOrgID is the Scaleway Organization ID, used for billing and resource ownership (legacy; prefer Project ID).
	envOrgID = "SCW_DEFAULT_ORGANIZATION_ID"

	// envAccessKey is the API access key for authenticating requests to Scaleway.
	envAccessKey = "SCW_ACCESS_KEY"

	// envSecretKey is the secret key associated with the access key, used for signing requests.
	envSecretKey = "SCW_SECRET_KEY"

	// envProjectID is the Scaleway Project ID, which groups resources and is the preferred way to organize infrastructure.
	envProjectID = "SCW_DEFAULT_PROJECT_ID"

	// envZone specifies the geographical region/zone where resources will be created (e.g., fr-par-1).
	envZone = "SCW_ZONE"

	// envBucket is a custom environment variable for specifying the name of an S3-compatible bucket.
	// This is not a standard Scaleway variable and is application-specific.
	envBucket = "SCW_BUCKET_NAME"
)

func main() {
	fmt.Println("moving snapshots to s3...")

	// Create a Scaleway client with credentials from environment variables.
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv(envOrgID)),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv(envAccessKey), os.Getenv(envSecretKey)),

		// Set the default project ID to organize resources under a specific project
		scw.WithDefaultProjectID(os.Getenv(envProjectID)),

		// Set the default zone where resources such as block volumes and snapshots are located
		scw.WithDefaultZone(scw.Zone(os.Getenv(envZone))),
	)
	if err != nil {
		panic(err)
	}

	instanceAPI := instance.NewAPI(client)

	snapList, err := instanceAPI.ListSnapshots(&instance.ListSnapshotsRequest{}, scw.WithAllPages())
	if err != nil {
		panic(err)
	}

	fmt.Printf("number of snapshots: %d\n", snapList.TotalCount)

	for _, snapshot := range snapList.Snapshots {
		fmt.Printf("snap %s\n", snapshot.Name)

		if snapshot.State == instance.SnapshotStateAvailable {
			fmt.Printf("Exporting snapshot %s (ID: %s) to bucket %s...\n", snapshot.Name, snapshot.ID, os.Getenv(envBucket))

			snap, err := instanceAPI.ExportSnapshot(&instance.ExportSnapshotRequest{
				SnapshotID: snapshot.ID,
				Bucket:     os.Getenv(envBucket),
				Key:        snapshot.Name + ".qcow2",
			})
			if err != nil {
				fmt.Printf("Failed to export snapshot %s: %v\n", snapshot.Name, err)

				continue
			}

			fmt.Printf("Successfully started export of %s to %s/%s\n", snap.Task.ID, os.Getenv(envBucket), snap.Task.Description)
		} else {
			fmt.Printf("Skipping snapshot %s (ID: %s) - status is %s, not available\n", snapshot.Name, snapshot.ID, snapshot.State.String())
		}
	}
}

// Check for mandatory variables before starting to work.
func init() {
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envZone, envProjectID, envBucket}

	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}
