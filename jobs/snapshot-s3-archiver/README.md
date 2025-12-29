# Scaleway Instance Snapshot Archiver

Automated Serverless Job to archive Scaleway Instance snapshots to Object Storage S3.

## Overview

This tool automatically finds available snapshots of Scaleway Instances volumes, exports them to a specified S3 bucket in `.qcow2` format, and deletes the original snapshot to optimize storage costs. It's designed to run as a Serverless Job on Scaleway and skips snapshots that have already been archived.

The main logic is implemented in `main.go`, which:
1. Loads configuration from environment variables.
2. Connects to Scaleway APIs using the Scaleway SDK.
3. Lists all available snapshots in the project.
4. Checks the target S3 bucket for already-archived snapshots.
5. Exports new snapshots to the bucket.
6. Deletes successfully exported snapshots to reduce storage costs.

## Features

- **Automated Export**: Finds available snapshots and exports them to an S3 bucket in `.qcow2` format.
- **Cost Optimization**: Deletes the source snapshot after successful export to reduce storage costs.
- **Idempotent**: Skips snapshots that are already archived in the bucket.
- **Serverless Ready**: Designed for [Scaleway Serverless Jobs](https://www.scaleway.com/en/serverless-jobs/).

## Step 1 : Build and push to Container registry

Serverless Jobs, like Serverless Containers (which are suited for HTTP applications), works
with containers. So first, use your terminal reach this folder and run the following commands:

```shell
# First command is to login to container registry, you can find it in Scaleway console
docker login rg.fr-par.scw.cloud/snapshot-s3-archiver -u nologin --password-stdin <<< "$SCW_SECRET_KEY"

# Here we build the image to push
docker buildx build --platform linux/amd64 -t rg.fr-par.scw.cloud/snapshot-s3-archiver/snapshot-s3-archiver:v1 .

# Push the image online to be used on Serverless Jobs
docker push rg.fr-par.scw.cloud/snapshot-s3-archiver/snapshot-s3-archiver:v1
```
> [!TIP]
> As we do not expose a web server and we do not require features such as auto-scaling, Serverless Jobs are perfect for this use case.
To check if everyting is ok, on the Scaleway Console you can verify if your tag is present in Container Registry.

## Step 2: Creating the Job Definition

On Scaleway Console on the following link you can create a new Job Definition: https://console.scaleway.com/serverless-jobs/jobs/create?region=fr-par

1. On Container image, select the image you created in the step before.
2. You can set the job definition name name to something clear.
3. Regarding the resources you can keep the default values, this job is fast and do not require specific compute power or memory.
4. To schedule your job for example every night at 2am, you can set the cron to `0 2 * * *`.
5. Important: advanced option, you need to set the following environment variables:

> [!TIP]
> For sensitive data like `SCW_ACCESS_KEY` and `SCW_SECRET_KEY` we recommend to inject them via Secret Manager, [more info here](https://www.scaleway.com/en/docs/serverless/jobs/how-to/reference-secret-in-job/).
| Variable | Description |
|---|---|
| `SCW_DEFAULT_ORGANIZATION_ID` | Organization ID . |
| `SCW_DEFAULT_PROJECT_ID` | Project ID (Recommended resource grouping). |
| `SCW_ACCESS_KEY` | IAM Access Key. |
| `SCW_SECRET_KEY` | IAM Secret Key. |
| `SCW_ZONE` | Zone of the snapshots (e.g., `fr-par-1`). |
| `SCW_BUCKET_NAME` | S3 Bucket name for archives. |
| `SCW_BUCKET_ENDPOINT` | S3 Endpoint (e.g., `s3.fr-par.scw.cloud`). |

* Then click "create job"

## Step 3: Run the job

On your created Job Definition, just click the button "Run Job" and within seconds it should be successful.

## Troubleshooting

If your Job Run state goes in error, you can use the "Logs" tab in Scaleway Console to get more informations about the error.

# Additional content

- [Jobs Documentation](https://www.scaleway.com/en/docs/serverless/jobs/how-to/create-job-from-scaleway-registry/)
- [Other methods to deploy Jobs](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/deploy-job/)
- [Secret key / access key doc](https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/)
- [CRON schedule help](https://www.scaleway.com/en/docs/serverless/jobs/reference-content/cron-schedules/)
-