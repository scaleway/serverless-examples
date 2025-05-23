# Memos application deployment using Terraform

[Memos](https://github.com/usememos/memos) is "an open-source, lightweight note-taking solution. The pain-less way to create your meaningful notes."

## Context

Memos requires a database and can be deployed on Serverless Containers. For cost optimised setup, for persistent storage this tutorial takes advantage of Scaleway Serverless SQL Database. Serverless Containers associated with Serverless Databases makes the setup cost efficient due to pay per use products.

## Requirements

- [Terraform installed](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli#install-terraform), more information about Scaleway Terraform Quickstart [here](https://www.scaleway.com/en/docs/terraform/quickstart/)
- Terraform requires configuration to access Scaleway ressource, [follow this reference documentation](https://www.scaleway.com/en/docs/terraform/reference-content/scaleway-configuration-file/)

## Usage

`main.tf` will:
* Create a new `memos` project
* Create required IAM permissions and key
* Create the Serverless SQL Database (postgres)
* Create the Serverless Container running `memos` configured to use the created SQL Database

Deploy the project using Terraform:

```bash
terraform init
terraform apply
```

## Delete changes

In case you need to delete everything created before (Database, Container, IAM configuration and Project), run:

```bash
terraform destroy
```
