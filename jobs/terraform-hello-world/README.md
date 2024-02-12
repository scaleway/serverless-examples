# Serverless Jobs Hello World with Terraform

This example demonstrates how to set up a Scaleway [Serverless Job](https://www.scaleway.com/en/serverless-jobs/) using [Terraform](https://www.terraform.io/).

It builds a custom image locally, pushes this to the Scaleway registry, then creates a job that runs this image on a schedule.

## Requirements

This example assumes you are familiar with how Serverless Jobs work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/jobs/quickstart/)

This example uses Terraform. Please set up your local environment as outlined in the docs for the [Scaleway Terraform Provider](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs).

You will also need a Scaleway API key, which can be configured using [Scaleway IAM](https://www.scaleway.com/en/docs/identity-and-access-management/iam/how-to/create-api-keys/). If you are using IAM policies, make sure the key has the permissions: `ServerlessJobsFullAccess`, `ContainerRegistryFullAccess`.

## Setup

Once your environment is set up, you need to export some environment variables to use your API key:

```console
export TF_VAR_access_key=<your api key>
export TF_VAR_secret_key=<your secret key>
export TF_VAR_project_id=<your project ID>
```

From there, you can run the following to set up your job:

```console
cd terraform

terraform init

terraform plan

terraform apply
```

You can then view your Job Definitions in the [Scaleway Console](https://console.scaleway.com/serverless-jobs/jobs).

The Job is set to run on a schedule, once every 5 minutes. Once the next schedule has run, you will see a "Job Run" listed for your Job Definition.

*Note* this job will keep running every 5 minutes, so you need to make sure you delete the job by running:

```console
cd terraform

terraform destroy
```
