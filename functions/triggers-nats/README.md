# Triggers NATS example

Simple example showing how to use NATS triggers with Scaleway Functions.

For complete examples of triggers in all function languages, see [triggers-getting-started](../triggers-getting-started).

The example function is triggered by a NATS queue, and will log the message body.

## Requirements

This example requires [Terraform](https://www.scaleway.com/en/docs/tutorials/terraform-quickstart/).

Also, NATS **must** be [activated](https://www.scaleway.com/en/docs/serverless/messaging/how-to/get-started/#how-to-create-a-nats-account) on your project.

## Setup

The Terraform configuration will create an example function, a NATS account, and a trigger. It will also write the NATS credentials to a file named `nats-creds`.

To authenticate with Scaleway, you can either set up the [Scaleway CLI](https://www.scaleway.com/en/cli/), from which Terraform can extract credentials, or you can export `SCW_ACCESS_KEY`, `SCW_SECRET_KEY` and `SCW_DEFAULT_PROJECT_ID`.

Once auth is set up, you can run:

```console
terraform init
terraform apply
```

You should be able to see your resources in the Scaleway console:

- NATS accounts can be found in the [MnQ section](https://console.scaleway.com/messaging/protocols/fr-par/sqs/queues)
- Functions can be found in the `triggers-nats` namespace in the [Serverless functions section](https://console.scaleway.com/functions/namespaces)

## Running

You can trigger your functions by sending messages to the associated NATS account. Below is a description of how to do this with our dummy `tests/send_messages.py` script.

### Setup

First you need to expose your NATS endpoint:

```console
export NATS_ENDPOINT=$(terraform output -raw nats_endpoint)
```

Then you can set up a Python environment in the `tests` directory, e.g.

```console
cd tests
python3 -m venv venv
source venv/bin/activate
pip3 install -r requirements.txt
```

### Sending messages

Now you can run the `send_messages.py` script to send a message to the NATS topic:

```console
python3 send_messages.py
```

### Viewing function logs

In your [Cockpit](https://console.scaleway.com/cockpit), you can access the logs from your queues and functions.

Navigate from your Cockpit to Grafana, and find the `Serverless Functions Logs` dashboard. There you should see logs from your function, printing the body of the NATS message.

## Cleanup

To delete all the resources used in this example, you can run the following from the root of the project:

```console
terraform destroy
```
