# Deploy an inference API as a container

## Step 1: Build and push API image to Scaleway's Registry

```bash
docker build -t rg.fr-par.scw.cloud/inference-api-images/inference-api:v1 .
docker login rg.fr-par.scw.cloud/inference-api-images -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
docker push rg.fr-par.scw.cloud/inference-api-images/inference-api:v1
```

## Step 2: Create and deploy a private inference container

Create `.tfvars` file in `/terraform` directory and put variable values in it:

```
region = "fr-par"
access_key = "<access-key>"
secret_key = "<secret_key>"
project_id = "<project_id>"
registry_image = "rg.fr-par.scw.cloud/registry-namespace-images/inference-api:v1"
```

Then perform:

```bash
cd terraform
terraform plan -var-file=testing.tfvars
terraform apply -var-file=testing.tfvars
```

## Test the inference API using HTTP calls

You can perform the following HTTP call:

```bash
curl -H "X-Auth-Token: $CONTAINER_TOKEN" -X POST "<scw_container_endpoint>" -H "Content-Type: application/json" -d '{"age": 44, "job": "blue-collar", "marital": "married", "education": "basic.4y", "default": "unknown", "housing": "yes", "loan": "no", "contact": "cellular", "month": "aug", "day_of_week": "thu", "duration": 210, "campaign": 1, "pdays": 999, "previous": "0", "poutcome": "nonexistent", "emp_var_rate": 1.4, "cons_price_idx": 93.444, "cons_conf_idx": -36.1, "euribor3m": 4.963, "nr_employed": 5228.1}'  
```