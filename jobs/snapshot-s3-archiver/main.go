package main

import (
	"context"
	"log/slog"
	"os"
	"slices"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const snapshotExtension = ".qcow2"

func main() {
	// Configure valid JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create Scaleway client using the implementation in config
	client, err := scw.NewClient(
		scw.WithDefaultOrganizationID(cfg.OrgID),
		scw.WithAuth(cfg.AccessKey, cfg.SecretKey),
		scw.WithDefaultProjectID(cfg.ProjectID),
		scw.WithDefaultZone(cfg.Zone),
	)
	if err != nil {
		slog.Error("Failed to create Scaleway client", "error", err)
		os.Exit(1)
	}

	slog.Info("Initializing instance API...")
	instanceAPI := instance.NewAPI(client)

	slog.Info("Reading all snapshots for the project...")
	snapList, err := instanceAPI.ListSnapshots(&instance.ListSnapshotsRequest{}, scw.WithAllPages())
	if err != nil {
		slog.Error("Failed to list snapshots", "error", err)
		os.Exit(1)
	}

	slog.Info("Reading all snapshots already in the bucket...")
	filesInBucket, err := listBucketFiles(cfg)
	if err != nil {
		slog.Error("Failed to list bucket files", "error", err)
		os.Exit(1)
	}

	processSnapshots(instanceAPI, cfg, snapList.Snapshots, filesInBucket)
}

func listBucketFiles(cfg *Config) ([]string, error) {
	minioClient, err := minio.New(cfg.BucketEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var files []string

	// List all objects in the bucket
	// Use a closed channel to signal cancellation if needed (not used here but good practice)
	// minioClient.ListObjects usually takes a channel for cancellation if needed, but here we pass context
	opts := minio.ListObjectsOptions{
		Recursive:    false,
		WithMetadata: true,
	}

	for object := range minioClient.ListObjects(ctx, cfg.BucketName, opts) {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, object.Key)
	}

	return files, nil
}

func processSnapshots(api *instance.API, cfg *Config, snapshots []*instance.Snapshot, filesInBucket []string) {
	for _, snapshot := range snapshots {
		logger := slog.With("snapshot_id", snapshot.ID, "snapshot_name", snapshot.Name)

		logger.Info("Checking snapshot")

		if snapshot.State != instance.SnapshotStateAvailable {
			logger.Info("Skipping snapshot (not available)", "status", snapshot.State.String())
			continue
		}

		filename := snapshot.Name + snapshotExtension

		if slices.Contains(filesInBucket, filename) {
			logger.Info("File already exists in bucket, deleting local snapshot")
			if err := api.DeleteSnapshot(&instance.DeleteSnapshotRequest{SnapshotID: snapshot.ID}); err != nil {
				logger.Error("Failed to delete snapshot", "error", err)
			}
			continue
		}

		logger.Info("Exporting snapshot to bucket")
		snap, err := api.ExportSnapshot(&instance.ExportSnapshotRequest{
			SnapshotID: snapshot.ID,
			Bucket:     cfg.BucketName,
			Key:        filename,
		})
		if err != nil {
			logger.Error("Failed to export snapshot", "error", err)
			continue
		}

		logger.Info("Successfully started export", "task_id", snap.Task.ID, "bucket", cfg.BucketName, "description", snap.Task.Description)
	}
}
