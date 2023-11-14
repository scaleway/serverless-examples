# Terraform, triggers and serverless containers

## Setup

Create a file `secrets.auto.tfvars` file holding your project ID, access key and secret key:

```
project_id = "your-project-id"
access_key = "your-access-key"
secret_key = "your-secret-key"
```

The API key you use only needs the following permission sets:

- `ContainerRegistryFullAccess`
- `ContainersFullAccess`
- `MessagingAndQueuingFullAccess`

**Optional**: you can also define in this file the language used for the container. By default, when not specified, it will use the `python` container.

```
container_language = "go"
```

The values for `container_language` can be either `python` or `go` (see inside [`docker/`](docker/) folder).

## Deploy

The deployment will do the following:

1. Create a Scaleway registry namespace
2. Build and deploy a container image with a Python/Go HTTP server
3. Deploy a public and private Serverless Container using the built image
4. Create Scaleway MnQ SQS queues and NATS subjects
5. Configure triggers from these SQS queues and NATS subjects to each container
6. Print the endpoints of each SQS queue, NATS subject and container

To run the deployment:

```shell
terraform init

terraform apply
```

## Running

You can test your triggers by sending messages to the SQS queues and NATS subjects configured with Terraform.

Below there is an example showing how to do this using a Python script in the [`tests/`](tests/) folder of this repo.

### Setup

First you need to export your queue credentials and URLs:

```shell
export AWS_ACCESS_KEY_ID=$(terraform output -raw sqs_admin_access_key)
export AWS_SECRET_ACCESS_KEY=$(terraform output -raw sqs_admin_secret_key)
export PUBLIC_QUEUE_URL=$(terraform output -raw public_queue)
export PRIVATE_QUEUE_URL=$(terraform output -raw private_queue)
export PUBLIC_SUBJECT=$(terraform output -raw public_subject)
export PRIVATE_SUBJECT=$(terraform output -raw private_subject)
export NATS_CREDS_FILE=$(terraform output -raw nats_creds_file)
```

You can then set up a Python environment in the `tests` directory:

```
cd tests
python3 -m venv venv
source venv/bin/activate
pip3 install -r requirements.txt
```

### Sending messages

You can now run the script to send messages to both queues with:

```
python3 send_messages.py
```

### Viewing function logs

In your [Cockpit](https://console.scaleway.com/cockpit), you can access the logs from your queues and functions.

Navigate from your Cockpit to Grafana, and find the `Serverless Containers Logs` dashboard. There you should see the "Hello World!" messages for the functions created by this example.

To get direct deep links to the function logs in Cockpit, you can also run the following from the root of the project:

```shell
echo "Go to $(terraform output --raw cockpit_logs_public_container) to see public container logs in cockpit"
echo "Go to $(terraform output --raw cockpit_logs_private_container) to see private container logs in cockpit"
```

Finally, you can modify `docker/python/server.py` (or `docker/go/server.go`) as you wish and run `terraform apply` to redeploy the new version of the container.

## Cleanup

Run `terraform destroy` to remove all resources created by this example.
