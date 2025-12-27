# Scaleway Instance Snapshot Archiver

Automated serverless job to archive Scaleway Instance snapshots to Object Storage S3.

## Features

- **Automated Export**: Finds available snapshots and exports them to an S3 bucket in `.qcow2` format.
- **Cost Optimization**: Deletes the source snapshot after successful export to reduce storage costs.
- **Idempotent**: Skips snapshots that are already archived in the bucket.
- **Serverless Ready**: Designed for [Scaleway Serverless Jobs](https://www.scaleway.com/en/serverless-jobs/).

## Configuration

Configure the job using environment variables:

| Variable | Description |
|---|---|
| `SCW_DEFAULT_ORGANIZATION_ID` | Organization ID (Legacy). |
| `SCW_DEFAULT_PROJECT_ID` | Project ID (Recommended resource grouping). |
| `SCW_ACCESS_KEY` | IAM Access Key. |
| `SCW_SECRET_KEY` | IAM Secret Key. |
| `SCW_ZONE` | Zone of the snapshots (e.g., `fr-par-1`). |
| `SCW_BUCKET_NAME` | S3 Bucket name for archives. |
| `SCW_BUCKET_ENDPOINT` | S3 Endpoint (e.g., `s3.fr-par.scw.cloud`). |

## Usage

### 1. Build

```bash
docker build -t snapshot-archiver .
```

### 2. Run Locally

Ensure all environment variables are set, then run:

```bash
go run .
```

### 3. Deploy

Push the image to your container registry and create a Serverless Job definition pointing to it with the required environment variables.