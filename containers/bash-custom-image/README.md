# Container bash script with custom image

This example builds a custom image, containing a Bash script, and deploys it as a Serverless Container.

The example installs the [MongoDB shell](https://www.mongodb.com/docs/mongodb-shell/), and runs a Bash script using `mongosh` to connect to a MongoDB instance.

This pattern can be copied and customise to install whatever tools and scripts you like.

## Requirements

- You have an account and are logged into the [Scaleway console](https://console.scaleway.com)
- You have created an API key in the [console](https://console.scaleway.com/iam/api-keys), with at least the `ContainerRegistryFullAccess`, and `ContainersFullAccess` permissions, plus access to the relevant project for Object Storage
- You have [Terraform](https://registry.terraform.io/providers/scaleway/scaleway/latest/docs) installed on your machine
- You have logged in to the Scaleway Container Registry (`scw registry login`)

If you are connecting to a MongoDB Atlas instance, you will also need to configure your [Network Access rules](https://www.mongodb.com/docs/atlas/security/ip-access-list/). You can either:

- Allowlist all the container IP ranges as described [in the FAQ](https://www.scaleway.com/en/docs/faq/serverless-containers/#can-i-whitelist-the-ips-of-my-containers)
- Open your instance to all traffic, adding the CIDR `0.0.0.0/0`

## Deploy on Scaleway

First you need to set the following environment variables:

```bash
# Variables for Scaleway project and API key
export TF_VAR_access_key=<your API access key>
export TF_VAR_secret_key=<your API secret key>
export TF_VAR_project_id=<your project id>

# Variables related to the MongoDB instance
export TF_VAR_mongo_hostname=<mongodb hostname>
export TF_VAR_mongo_username=<mongodb username>
export TF_VAR_mongo_password=<mongodb password>
```

Deployment can be done by running:

```bash
cd terraform

terraform init

terraform plan

terraform apply
```

You can then query your function by running:

```bash
curl $(terraform output -raw endpoint)
```

This will print the stats about your MongoDB instance.
