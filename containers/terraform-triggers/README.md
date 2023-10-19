# Terraform + SQS + NATS triggers

## Setup

Create a file `secrets.auto.tfvars` file in this directory with:

```
project_id = "your-project-id"
access_key = "your-access-key"
secret_key = "your-secret-key"
```

## Deploy

The deployment will do the following:

1. Create a Scaleway registry namespace
2. Build and deploy a container image with a Python HTTP server
3. Create Scaleway MnQ SQS queues
4. Deploy a public and private Serverless Container
5. Configure triggers to each container
6. Print the endpoints of each queue and each container

To run the deploy:

```
scw registry login

terraform init

terraform plan
terraform apply
```
