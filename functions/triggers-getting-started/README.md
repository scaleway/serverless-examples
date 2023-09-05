# Triggers Getting Started

Simple starter examples that showcase using SQS triggers in all Scaleway Functions supported languages.

In each example, a function is triggered by a SQS queue with a message containing a number. The function will then print the factorial of this number to the logs.

## Requirements

This example requires [Terraform](https://www.scaleway.com/en/docs/tutorials/terraform-quickstart/).

## Setup

The Terraform configuration will deploy a function for each language, showing how to use triggers with each language.

It will also create a SQS queue per function to trigger it.

```console
terraform init
terraform apply
```

You should be able to see your functions in the Scaleway console.

## Running

You can use the `tests/send_messages.py` script to send a message to the SQS queue of each function.

```console
export AWS_ACCESS_KEY_ID=$(terraform output -raw sqs_access_key)
export AWS_SECRET_ACCESS_KEY=$(terraform output -raw sqs_secret_key)
python tests/send_messages.py
```

In Cockpit, you should see the functions being triggered and the result of the factorial being printed in the logs.

```console
```

## Cleanup

```console
terraform destroy
```
