# Terraform NGINX hello world

This demonstrates a simple example of deploying a container with Terraform on [Scaleway Serverless Containers](https://www.scaleway.com/en/serverless-containers/).

## Requirements

This example assumes you are familiar with how serverless containers work. If needed, you can
check the [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/containers/quickstart/).

This example uses Terraform. Please install and configure Terraform before trying out the example.

## Setup

Set up your scaleway credentials with:
```sh
export TF_VAR_access_key=<scw-access-key>
export TF_VAR_secret_key=<scw-secret-key>
export TF_VAR_project_id=<scw-project-id>
```

Set up Terraform with:
```sh
terraform init
```

## Deployment

Once your environment is set up, you can:

Run plan, have a look at what will be created:
```sh
terraform plan
```

Deploy with:
```sh
terraform apply
```

When the deployment is complete, you can check the deployment succeeded either by:

i) curl the container's endpoint with:
```sh
curl $(terraform output -raw endpoint)
```
ii) hit it from a browser and see the NGINX default page.
