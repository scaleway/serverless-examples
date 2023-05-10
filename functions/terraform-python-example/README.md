# Terraform Python example

In this tutorial you will discover an example of Instance automation using Python and Terraform:

* Automatically shut down / start instances

## Requirements

* You have an account and are logged into the [Scaleway console](https://console.scaleway.com)
* You have [generated an API key](/console/my-project/how-to/generate-api-key/)
* You have [Python](https://www.python.org/) installed on your machine
* You have [Terraform](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs) installed on your machine
* You are familiar with Instance API that can be found in the [developers documentation](https://developers.scaleway.com/en/products/instance/api/#get-2c1c6f).

## Context

In this tutorial, we will simulate a project with a production environment that will be running all the time and a development environment that will be turn off on week-ends to save costs.

## Project Structure

1. This folder stores your configuration as explained in the terraform documentation.

2. This folder contains 5 files to configure your infrastructure:
  a. 'main.tf': will contain the main set of configurations for your project. Here, it will be our instance.
  b. 'provider.tf': Terraform relies on plugins called “providers” to interact with remote systems.
  c. 'variables.tf': will contain the variable definitions for your project. Since all Terraform values must be defined, any variables that are not given a default value will become required arguments.
  d. 'terraform.tfvars': allows you to set the actual value of the variables.
3. Create the following folder:
  a. 'function': to store your function code.
  b. 'files': (generated in `main.tf` when creating function's zip file) to temporarily store your zip function code.

## Deploy your infrastructure

Now that everything is set up, deploy everything using Terraform

1. Add your Scaleway credentials to your environment variables

```bash
export SCW_ACCESS_KEY="<your-secret-key>"
export SCW_SECRET_KEY="<your-access-key>"
```

2. Initialize Terraform:

```bash
terraform init
```

3. Let terraform verify your configuration:

```bash
terraform plan
```

4. Deploy your infrastructure:

```bash
terraform apply
````

5. In order to remove your infrastructure resources, you can use:

```bash
terraform destroy
```