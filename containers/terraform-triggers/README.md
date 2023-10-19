# Terraform, triggers and serverless containers

## Setup

Create a file `secrets.auto.tfvars` file holding your project ID, access key and secret key:

```
project_id = "your-project-id"
access_key = "your-access-key"
secret_key = "your-secret-key"
```

## Deploy

The deployment will do the following:

1. Create a Scaleway registry namespace
2. Build and deploy a container image with a Python HTTP server
3. Deploy a public and private Serverless Container using the built image
4. Create Scaleway MnQ SQS queues
5. Configure triggers from these queues to each container
6. Print the endpoints of each queue and each container

To run the deploy:

```
terraform init

terraform plan
terraform apply
```
