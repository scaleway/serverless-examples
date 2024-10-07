# Node Terraform

A simple example of deploying a Node Serverless function using Terraform. The function is a simple `rss` feed that filters the content of a source feed and returns the filtered content.

## Requirements

- [Node.js](https://nodejs.org/en/download/)
- [Terraform](https://learn.hashicorp.com/terraform/getting-started/install.html)
- A configured Scaleway Profile. You can find more information [here](https://www.scaleway.com/en/docs/developer-tools/terraform/reference-content/scaleway-configuration-file/#how-to-set-up-the-configuration-file)

## Usage

Run the function locally:

```bash
cd function
npm install --include=dev
npm start
```

In another terminal, you can the following command to test the function:

```bash
curl http://localhost:8081
```

Deploy the function using Terraform:

```bash
terraform init
terraform apply
```

The function is a `rss` feed that can be accessed via a RSS reader. The URL of the feed is displayed in the output of the Terraform apply command.

## Cleanup

```bash
terraform destroy
```
