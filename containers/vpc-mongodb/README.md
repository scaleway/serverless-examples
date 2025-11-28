# Serverless Containers with Scaleway Managed MongoDB®

A simple Deno application demonstrating how to connect to a Scaleway Managed MongoDB® instance using Mongoose. It relies on Scaleway's Private Networks for secure communication.

## Deploying

### Prerequisites

This Terraform pushes a Docker image to a Scaleway Container Registry before deploying the Serverless Container.

Make sure to log in to the Scaleway Container Registry beforehand:

```bash
scw registry login
```

### Deploy with Terraform

```bash
terraform init
terraform apply
```

That's it! You should be able to access the deployed Serverless Container's endpoint from the Scaleway Console.

## Using the Application

Once deployed, you can interact with the application using HTTP requests.

```bash
CONTAINER_URL="https://your-container-url"
# Check MongoDB® connection
curl $CONTAINER_URL/check_connection
# Add some people
curl $CONTAINER_URL/person/george
curl $CONTAINER_URL/person/alice
# List all people
curl $CONTAINER_URL/people
```

## Cleaning Up

Make sure to destroy the Terraform-managed infrastructure when you're done:

```bash
terraform destroy
```
