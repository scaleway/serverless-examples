# Access S3 from a PHP function

This example demonstrates how we can access S3 from PHP functions.

We will:

1. Set up an S3 bucket using Terraform
2. Deploy a PHP function using Terraform
3. Invoke the PHP function which writes to the S3 bucket

This example assumes you are familiar with how serverless functions work. If needed, you can check the [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

## Requirements

- [Scaleway CLI](https://github.com/scaleway/scaleway-cli). See the [install instructions](https://github.com/scaleway/scaleway-cli#installation)
- [Terraform](https://www.terraform.io/). See the [Terraform Scaleway provider docs](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs)
- `make` to run targets in the `Makefile`

## Setup

First you need to set up your Terraform environment with your Scaleway credentials by exporting some environment variables:

```
export TF_VAR_project_id=<your project id>
export TF_VAR_access_key=<your access key>
export TF_VAR_secret_key=<your secret key>
```

You can then set up, plan and apply the Terraform configuration which will create an S3 bucket, as well as the PHP function:

```
# Initialise Terraform
cd terraform
terraform init

# Check what changes Terraform will make
terraform plan

# Apply the changes
terraform apply
```

Then you can wait for the function to be deployed, and run:

```
curl https://$(terraform output -raw function_url)
```

You can then check [your `php-s3-example` bucket](https://console.scaleway.com/object-storage/buckets/fr-par/php-s3-example/explorer) to see the key that was written by your function.

You will also see the S3 response printed in the function logs.
