# Scaleway Instance Snapshot Backup to S3

This project exports available Scaleway Instance snapshots to an S3-compatible bucket (e.g., Scaleway Object Storage), and optionally deletes the snapshot afterward if it's already backed up. It's designed to run as a **Scaleway Serverless Job**, making it ideal for automated, scheduled backups.

---

## 📦 Features

- Lists all available block storage snapshots in a project.
- Checks if a snapshot with the same name already exists in the target bucket.
- Exports missing snapshots to the bucket in `.qcow2` format.
- Deletes local snapshot after successful export (if not already in bucket).
- Uses environment variables for full configuration.
- Built to run in a container on [Scaleway Serverless Jobs](https://www.scaleway.com/en/serverless-jobs/).

---

## ⚙️ Environment Variables

You must set the following environment variables when deploying the job:

| Variable | Description |
|--------|-------------|
| `SCW_DEFAULT_ORGANIZATION_ID` | Your Scaleway Organization ID (legacy; prefer project ID). |
| `SCW_DEFAULT_PROJECT_ID` | Your Scaleway Project ID (preferred way to group resources). |
| `SCW_ACCESS_KEY` | API access key (from IAM). |
| `SCW_SECRET_KEY` | API secret key (from IAM). |
| `SCW_ZONE` | Zone where your snapshots are located (e.g., `fr-par-1`). |
| `SCW_BUCKET_NAME` | Name of the S3 bucket to store exported snapshots. |
| `SCW_BUCKET_ENDPOINT` | S3 endpoint (e.g., `s3.fr-par.scw.cloud`). |

> 🔐 **Security Tip**: Use IAM API keys with minimal required permissions.

---

## 🛠️ Build & Deploy to Scaleway Serverless Jobs

### 1. Build the Docker Image

```bash
docker build -t snapshot-s3-backup .
```

### 2. Tag and Push to Scaleway Container Registry (or any registry)

```bash
# Example using Scaleway CR
docker tag snapshot-s3-backup fr-par.scw.cloud/your-registry/snapshot-s3-backup:v1
docker push fr-par.scw.cloud/your-registry/snapshot-s3-backup:v1
```

> Replace `your-registry` with your actual container registry name.

### 3. Create the Serverless Job

Use the Scaleway CLI or Console:

#### Using `scw` CLI:

```bash
scw job create \
  name=backup-snapshots \
  image=fr-par.scw.cloud/your-registry/snapshot-s3-backup:v1 \
  memory-limit=512Mi \
  cpu-limit=500m \
  environment='{
    "SCW_DEFAULT_PROJECT_ID": "your-project-id",
    "SCW_ACCESS_KEY": "your-access-key",
    "SCW_SECRET_KEY": "your-secret-key",
    "SCW_ZONE": "fr-par-1",
    "SCW_BUCKET_NAME": "my-backup-bucket",
    "SCW_BUCKET_ENDPOINT": "s3.fr-par.scw.cloud"
  }'
```

### 4. (Optional) Schedule the Job

Schedule it to run daily using a cron trigger:

```bash
scw scheduler trigger create-cron \
  job-id=your-job-id \
  schedule="0 2 * * *" \
  name=daily-snapshot-backup
```

This runs the job every day at 2 AM.

---

## 📁 Output Format

Each snapshot is exported as:
```
<snapshot-name>.qcow2
```

Example:
```
my-server-disk-2025-04-05.qcow2
```

---

## ✅ Example Use Case

Run nightly to:
1. Export new snapshots to object storage.
2. Clean up old snapshots once safely backed up.
3. Reduce storage costs and improve disaster recovery.

---

## 🧪 Local Testing (Optional)

Set environment variables:

```bash
export SCW_DEFAULT_PROJECT_ID=...
export SCW_ACCESS_KEY=...
export SCW_SECRET_KEY=...
export SCW_ZONE=fr-par-1
export SCW_BUCKET_NAME=my-backup-bucket
export SCW_BUCKET_ENDPOINT=s3.fr-par.scw.cloud
```

Run:

```bash
go run main.go
```

---