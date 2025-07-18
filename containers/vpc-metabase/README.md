# Metabase on VPC

This example shows how to deploy a Metabase instance on Serverless Containers that connects to a PostgreSQL database running in a private network. The setup uses Terraform to manage the infrastructure.

## Prerequisites

- A Scaleway account
- [Terraform](https://www.terraform.io/downloads.html) installed

## Setup

Deploy the PostgreSQL database and private network:

```bash
cd containers/vpc-metabase
terraform init
terraform apply
```

After, the deployment, you can find the Metabase URL in the output:

```bash
terraform output metabase_container_url
```

That's it! You can now access your Metabase instance at the provided URL 🎉. Please refer to the official [Metabase documentation](https://www.metabase.com/docs/latest/) for more information.
