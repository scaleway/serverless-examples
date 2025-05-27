# Memos application deployment using Terraform

[Memos](https://github.com/usememos/memos) is "an open-source, lightweight note-taking solution. The pain-less way to create your meaningful notes."

## Context

Memos requires a database and can be deployed on Serverless Containers. For a cost-optimized setup with persistent storage, this tutorial takes advantage of Scaleway's Serverless SQL Database. Serverless Containers combined with Serverless Databases make the setup cost-efficient thanks to their pay-per-use model.

## Requirements

- [Terraform installed](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli#install-terraform), more information about Scaleway Terraform Quickstart [here](https://www.scaleway.com/en/docs/terraform/quickstart/)
- Terraform requires configuration to access Scaleway ressource, [follow this reference documentation](https://www.scaleway.com/en/docs/terraform/reference-content/scaleway-configuration-file/)

## Usage

`main.tf` will:
* Create a new `memos` project
* Create required IAM permissions and key
* Create the Serverless SQL Database (postgres)
* Create the Serverless Container running `memos` and configured to use the newly created SQL Database

Deploy the project using Terraform:

```bash
terraform init
terraform apply
```

## Delete changes

In case you need to delete everything created earlier (Database, Container, IAM configuration and Project), run the following command:

```bash
terraform destroy
```
