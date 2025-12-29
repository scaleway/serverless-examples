package main

import (
	"context"
	"log/slog"
	"os"
	"slices"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/scaleway/scaleway-sdk-go/api/block/v1"
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

	slog.Info("Initializing block API...")
	blockAPI := block.NewAPI(client)

	slog.Info("Reading all snapshots for the project...")
	snapList, err := blockAPI.ListSnapshots(&block.ListSnapshotsRequest{
		Zone:      cfg.Zone,
		ProjectID: &cfg.ProjectID,
	}, scw.WithAllPages())
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

	processSnapshots(blockAPI, cfg, snapList.Snapshots, filesInBucket)
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

func processSnapshots(api *block.API, cfg *Config, snapshots []*block.Snapshot, filesInBucket []string) {
	for _, snapshot := range snapshots {
		logger := slog.With("snapshot_id", snapshot.ID, "snapshot_name", snapshot.Name)

		logger.Info("Checking snapshot")

		if snapshot.Status != block.SnapshotStatusAvailable {
			logger.Info("Skipping snapshot (not available)", "status", snapshot.Status.String())
			continue
		}

		filename := snapshot.Name + snapshotExtension

		if slices.Contains(filesInBucket, filename) {
			logger.Info("File already exists in bucket, deleting local snapshot")
			err := api.DeleteSnapshot(&block.DeleteSnapshotRequest{
				SnapshotID: snapshot.ID,
				Zone:       snapshot.Zone,
			})
			if err != nil {
				logger.Error("Failed to delete snapshot", "error", err)
			}
			continue
		}

		logger.Info("Exporting snapshot to bucket")
		snap, err := api.ExportSnapshotToObjectStorage(&block.ExportSnapshotToObjectStorageRequest{
			SnapshotID: snapshot.ID,
			Bucket:     cfg.BucketName,
			Key:        filename,
			Zone:       snapshot.Zone,
		})
		if err != nil {
			logger.Error("Failed to export snapshot", "error", err)
			continue
		}

		logger.Info("Successfully started export", "task_id", snap.ID, "bucket", cfg.BucketName)
	}
}
