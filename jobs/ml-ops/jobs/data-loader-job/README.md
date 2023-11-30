# Data Loader Job

## Before pushing docker image to private registry

Create and fill `.env` file with these variables with appropriate values:

```bash
SCW_ACCESS_KEY=my_access_key
SCW_SECRET_KEY=my_secret_key
```

## Define and run data loader job on the console

Use these environment variables for your job:

```text
SCW_S3_BUCKET=<data-store-name>
SCW_REGION=fr-par
SOURCE_FILE_NAME=bank_telemarketing.csv
```