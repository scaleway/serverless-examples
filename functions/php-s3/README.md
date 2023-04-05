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

To configure Terraform to use your project, edit `terraform/vars/main.tfvars` to set your project ID:

```
# Replace with your project ID
project_id = "12ef4x91-yh12-1234-g22g-83er2q4z51ec"
```

You then need to export your Scaleway access key and secret key:

```
export TF_VAR_access_key=<your access key>
export TF_VAR_secret_key=<your secret key>
```

You can then set up, plan and apply the Terraform configuration which will create an S3 bucket, as well as the PHP function:

```
# Initialise Terraform
make tf-init

# Check what changes Terraform will make
make tf-plan

# Apply the changes
make tf-apply
```

This will also output a script, `curl.sh` which you can run to call your function:

```
./curl.sh
```

You can then check [your `php-s3-example` bucket](https://console.scaleway.com/object-storage/buckets/fr-par/php-s3-example/explorer) to see the key that was written by your function.

You will also see the S3 response printed in the function logs.
