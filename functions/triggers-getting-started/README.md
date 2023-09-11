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
...
DEBUG - php: factorial of 17 is 355687428096000 source=user stream=stdout  
2023-09-11 10:36:19.994 DEBUG - Function Triggered: / source=core
2023-09-11 10:36:19.993 DEBUG - php: factorial of 13 is 6227020800 source=user stream=stdout
2023-09-11 10:36:19.991 DEBUG - Function Triggered: / source=core
2023-09-11 10:36:19.977 DEBUG - php: factorial of 12 is 479001600 source=user stream=stdout
2023-09-11 10:36:19.976 DEBUG - php: factorial of 11 is 39916800 source=user stream=stdout
2023-09-11 10:36:19.975 DEBUG - php: factorial of 10 is 3628800 source=user stream=stdout
2023-09-11 10:36:19.964 DEBUG - Function Triggered: / source=core
2023-09-11 10:36:19.954 DEBUG - php: factorial of 3 is 6 source=user stream=stdout
2023-09-11 10:36:19.954 DEBUG - php: factorial of 4 is 24 source=user stream=stdout
2023-09-11 10:36:19.948 DEBUG - php: factorial of 0 is 1 source=user stream=stdout
... (truncated)
```

Configuring and managing triggers is free, you only pay for the execution of your functions. For more information, please consult the [Scaleway Serverless pricing](https://www.scaleway.com/en/pricing/?tags=serverless). You will also be billed for the SQS queue usage when sending messages to it.

## Cleanup

```console
terraform destroy
```
