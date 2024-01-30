# Container used to upload files to S3

This container does the following:

* Read a file from an HTTP request form
* Store the file in S3

## Requirements

- You have an account and are logged into the [Scaleway console](https://console.scaleway.com)
- You have created an API key in the [console](https://console.scaleway.com/iam/api-keys), with at least the `ObjectStorageFullAccess`, `ContainerRegistryFullAccess`, and `FunctionsFullAccess` permissions, plus access to the relevant project for Object Storage
- You have [Terraform](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs) installed on your machine
- You have logged in to the Scaleway Container Registry (`scw registry login`)

## Deploy on Scaleway

First you need to set the following environment variables:

```bash
export TF_VAR_access_key=<your API access key>
export TF_VAR_secret_key=<your API secret key>
export TF_VAR_project_id=<your project id>
```

Deployment can be done by running:

```bash
terraform init

terraform plan

terraform apply
```

You can then query your function by running:

```bash
# Upload the requirements file
curl -F file=@requirements.txt $(terraform output -raw function_url)
```

You can get the bucket name with:

```bash
terraform output -raw bucket_name
```

You should then see the `requirements.txt` file uploaded to your bucket.
