# Terraform, triggers and serverless containers

## Setup

Create a file `secrets.auto.tfvars` file holding your project ID, access key and secret key:

```
project_id = "your-project-id"
access_key = "your-access-key"
secret_key = "your-secret-key"
```

## Deploy

The deployment will do the following:

1. Create a Scaleway registry namespace
2. Build and deploy a container image with a Python HTTP server
3. Deploy a public and private Serverless Container using the built image
4. Create Scaleway MnQ SQS queues
5. Configure triggers from these queues to each container
6. Print the endpoints of each queue and each container

To run the deploy:

```
terraform init

terraform apply
```

## Tests

```shell
export AWS_ACCESS_KEY_ID=$(terraform output -raw sqs_admin_access_key)
export AWS_SECRET_ACCESS_KEY=$(terraform output -raw sqs_admin_secret_key)
export PUBLIC_QUEUE_URL=$(terraform output -raw public-queue)
export PRIVATE_QUEUE_URL=$(terraform output -raw private-queue)

cd tests/
python3 -m venv env
source env/bin/activate
python3 -m pip install -r requirements.txt
python3 send_messages.py
```

You can then check logs of your containers in your cockpit (in project defined in `secrets.auto.tfvars`). You should see the "Hello World!" messages there!

From there, you can modify `docker/server.py` as you want and run `terraform apply` to deploy the modified container.

## Cleanup

Run `terraform destroy` to remove all resources created by this example.
