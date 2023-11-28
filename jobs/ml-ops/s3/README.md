# Object Storage buckets 

## Step 1: Create object storage resources with Terraform

Create `.tfvars` file in `/terraform` directory and put variable values in it. Then perform:

```bash
cd terraform
terraform plan -var-file="file_name.tfvars"
terraform apply -var-file="file_name.tfvars"
```

This will create three buckets:

* `data-store`: used to store data files.
* `model-registry`: used to store  a dump of trained ML model.
* `performance-monitoring`: used to store performance of the trained ML model on test (unobserved) data.

## Step 2: Push local data to data store bucket

Create and fill `.env` file with these variables with appropriate values:

```bash
SCW_ACCESS_KEY=my_access_key
SCW_SECRET_KEY=my_secret_key
SCW_S3_BUCKET=data-store
SCW_REGION=fr-par
SOURCE_FILE_NAME=bank_telemarketing.csv
```

Push the data file to the created data store bucket using:

```bash
cd terraform
python main.py
```
