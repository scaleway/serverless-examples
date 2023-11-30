# Machine Learning job for binary classification use case

## Before pushing docker image to private registry

Create and fill `.env` file with these variables with appropriate values:

```bash
SCW_ACCESS_KEY=my_access_key
SCW_SECRET_KEY=my_secret_key
```

## Define and run an ML job

You can create a job definition on the console using the private registry image. Run the job and check that training artifacts are uploaded to object storage buckets.

Define these environment variables during job run:

```text
SCW_S3_BUCKET_DATA=<data-store-name>
DATA_FILE_NAME=bank_telemarketing.csv
SCW_S3_BUCKET_MODEL=<model-registry-name>
MODEL_FILE_NAME=classifier.pkl
SCW_S3_BUCKET_PERF=<perf-monitor-name>
SCW_REGION=fr-par
```