# Machine Learning job for binary classification use case

## Step 1: Build and push ML training image to Scaleway's Registry

Create an fill a `.env` file as it follows:

```bash
SCW_ACCESS_KEY=<scw_access_key>
SCW_SECRET_KEY=<scw_secret_key>
SCW_S3_BUCKET_DATA=data-store
DATA_FILE_NAME=bank_telemarketing.csv
SCW_S3_BUCKET_MODEL=model-registry
MODEL_FILE_NAME=classifier.pkl
SCW_S3_BUCKET_PERF=performance-monitoring
SCW_REGION=fr-par
```

Then build and push job image to registry:

```bash
docker build -t rg.fr-par.scw.cloud/ml-job-images/ml-job:v1 .
docker login rg.fr-par.scw.cloud/ml-job-images -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
docker push rg.fr-par.scw.cloud/ml-job-images/ml-job:v1
```

## Step 2: Define and run an ML job

You can create a job definition on the console using the private registry image. Run the job and check that training artifacts are uploaded to object storage buckets.
