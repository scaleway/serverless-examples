# Scaleway private registry for docker images

## Inference API image registry

The registry is used to store inference API image. This image can be deployed a serverless containers.

## ML model training

The registry is used to store image to train a machine learning model. This image can be deployed as a serverless job.

## Create registry namespaces

Create `.tfvars` file in `/terraform` directory and put variable values in it:

```
region = "fr-par"
access_key = "<access-key>"
secret_key = "<secret_key>"
project_id = "<project_id>"
```

Then perform:

```bash
cd terraform
terraform plan -var-file=testing.tfvars
terraform apply -var-file=testing.tfvars
```

