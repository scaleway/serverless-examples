package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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

	envBucketEndpoint = "SCW_BUCKET_ENDPOINT"
)

func main() {
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

	fmt.Println("Initializing instance API...")

	instanceAPI := instance.NewAPI(client)

	fmt.Println("Reading all snapshots for the project...")

	snapList, err := instanceAPI.ListSnapshots(&instance.ListSnapshotsRequest{}, scw.WithAllPages())
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading all snapshots already in the bucket...")

	filesInBucket, err := listBucketFiles()
	if err != nil {
		panic(err)
	}

	const snapshotExtension = ".qcow2"

	for _, snapshot := range snapList.Snapshots {
		fmt.Printf("Checking for snapshot %s\n", snapshot.Name)

		if snapshot.State == instance.SnapshotStateAvailable {
			// Check if file already exists in bucket
			if slices.Contains(filesInBucket, snapshot.Name+snapshotExtension) {
				fmt.Printf("File %s already exists in bucket, can delete the snapshot and skip it\n", snapshot.Name+snapshotExtension)

				err = instanceAPI.DeleteSnapshot(&instance.DeleteSnapshotRequest{
					SnapshotID: snapshot.ID,
				})
				if err != nil {
					panic(err)
				}

				continue
			}

			fmt.Printf("File %s not present in the  bucket, expording it to the bucket...\n", snapshot.Name+".qcow2")

			snap, err := instanceAPI.ExportSnapshot(&instance.ExportSnapshotRequest{
				SnapshotID: snapshot.ID,
				Bucket:     os.Getenv(envBucket),
				Key:        snapshot.Name + snapshotExtension,
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
	mandatoryVariables := [...]string{
		envOrgID,
		envAccessKey,
		envSecretKey,
		envZone,
		envProjectID,
		envBucket,
		envBucketEndpoint,
	}

	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}

func listBucketFiles() ([]string, error) {
	// Retrieve S3-compatible endpoint and credentials from environment
	endpoint := os.Getenv(envBucketEndpoint)
	accessKeyID := os.Getenv(envAccessKey)
	secretAccessKey := os.Getenv(envSecretKey)

	// Create new MinIO client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	// Set up context and result slice
	ctx := context.Background()
	var files []string

	// Channel to signal listing completion
	doneCh := make(chan struct{})
	defer close(doneCh)

	// List all objects in the bucket
	for object := range minioClient.ListObjects(ctx, os.Getenv(envBucket), minio.ListObjectsOptions{
		Recursive:    false,
		WithMetadata: true,
	}) {
		if object.Err != nil {
			return nil, object.Err
		}

		files = append(files, object.Key)
	}

	return files, nil
}
