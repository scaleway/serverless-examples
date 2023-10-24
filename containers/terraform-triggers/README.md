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
4. Create Scaleway MnQ SQS queues
5. Configure triggers from these queues to each container
6. Print the endpoints of each queue and each container

To run the deployment:

```shell
terraform init

terraform apply
```

## Tests

You can send some messages by using the provided python script in [`tests/`](tests/) folder.

```shell
export AWS_ACCESS_KEY_ID=$(terraform output -raw sqs_admin_access_key)
export AWS_SECRET_ACCESS_KEY=$(terraform output -raw sqs_admin_secret_key)
export PUBLIC_QUEUE_URL=$(terraform output -raw public_queue)
export PRIVATE_QUEUE_URL=$(terraform output -raw private_queue)

python3 -m venv tests/env
source tests/env/bin/activate
python3 -m pip install -r tests/requirements.txt
python3 tests/send_messages.py
```

You can then check logs of your containers in your cockpit (assuming you have activated it in your project):

```shell
echo "Go to $(terraform output --raw cockpit_logs_public_container) to see public container logs in cockpit"
echo "Go to $(terraform output --raw cockpit_logs_private_container) to see private container logs in cockpit"
```

You should be able to see the "Hello World!" messages there.

Finally, you can modify `docker/python/server.py` (or `docker/go/server.go`) as you wish and run `terraform apply` to redeploy the new version of the container.

## Cleanup

Run `terraform destroy` to remove all resources created by this example.
