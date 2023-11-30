# Machine Learning job for binary classification use case

## Define and run an ML job

You can create a job definition on the console using the private registry image. Run the job and check that training artifacts are uploaded to object storage buckets.

Use these job environment variables:

```bash
SCW_S3_BUCKET_DATA=data-store
SCW_S3_BUCKET_MODEL=model-registry
SCW_S3_BUCKET_PERF=performance-monitoring
```