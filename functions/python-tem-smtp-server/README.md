# Send Transactional Emails from a Serverless Function

This example demonstrates how to send TEM with an SMTP server from Python functions.

It assumes that you are familiar with how Serverless Functions and Scaleway Transactional Email work. 
If needed, you can check the Scaleway official documentation for serverless functions [here](https://www.scaleway.com/en/docs/serverless/functions/quickstart/) 
and for TEM [here](https://www.scaleway.com/en/docs/managed-services/transactional-email/quickstart/).

## Requirements

* You have generated an API key with the permission `TransactionalEmailFullAccess`
* You have [Python](https://www.python.org/) installed on your machine
* You have [Terraform](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs) installed on your machine

## Setup

You have to configure your domain with Transactional Email (tutorial available [here](https://www.scaleway.com/en/docs/managed-services/transactional-email/quickstart/))

Then, edit the file `handler.py` to set the sender and the recipient of the email. 

Also, depending on your SMTP server, you might also need to change the value of the variables `host` and `port`.

## Testing with serverless offline for Python

In order to test your function locally before the deployment, you can install our offline testing library with:

```bash
pip install scaleway_functions_python==0.2.0
```

Export your environment variables and then launch your function locally:
```bash
export TEM_PROJECT_ID=<the Project ID of the Project in which the TEM domain was created>
export SECRET_KEY=<the secret key of the API key of the project used to manage your TEM domain>

python handler.py
```

Test your local function using `curl`:

```bash
curl http://localhost:8080
```

This should email the recipient defined in `handler.py`.

## Deploy

Use the Terraform configuration to deploy the function.

```shell
terraform init
terraform apply -var "scw_project_id=$TEM_PROJECT_ID" -var "scw_secret_key=$SECRET_KEY"
```

## Call the function

When the deployment is complete, you can check the deployment succeeded either by:

i) curl the container's endpoint with:
```sh
curl $(terraform output -raw endpoint)
```
ii) hit it from a browser.

Doing so will send an e-mail to the recipient defined in the file `handler.py`.
