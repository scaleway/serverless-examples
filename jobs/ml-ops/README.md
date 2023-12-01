# Serverless MLOps

In this example, we train and deploy a binary classification inference model using Scaleway Serverless. To do this, we use the following resources:

1. Serverless Job for training
2. Serverless Job to populate data in S3
3. Serverless Container for inference

We use object storage to share data between the two.

## Context

In this example we use a bank telemarketing dataset to predict if a client would engage in a term deposit subscription.

This dataset records marketing phone calls made to clients. The outcome of the phone call is in shown in the `y` column:

* `0` : no subscription
* `1` : subscription

## Data Source

The dataset has many versions and is open-sourced and published [here](http://archive.ics.uci.edu/dataset/222/bank+marketing) on the UCI Machine Leaning repository and is close to the one analyzed in the following research work:

* [Moro et al., 2014] S. Moro, P. Cortez and P. Rita. A Data-Driven Approach to Predict the Success of Bank Telemarketing. Decision Support Systems, Elsevier, 62:22-31, June 2014

## Running the example

### Step 1. Provision resources with Terraform

Set your Scaleway access key, secret key and project ID in environment variables:

```console
export TF_VAR_access_key=<your-access-key>
export TF_VAR_secret_key=<your-secret-key>
export TF_VAR_project_id=<your-project-id>

cd terraform
terraform init
terraform plan
terraform apply
```

### Step 2. Run the data and training Jobs

*At the time of writing, the Scaleway CLI does not support Jobs, so we use a Python script*

```
cd scripts

python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

python3 run upload
python3 run training
```

You can then check your Job runs in the [Jobs Console](https://console.scaleway.com/serverless-jobs/jobs).

### Step 4. Use the inference API

```
export INFERENCE_URL=$(terraform output endpoint)

curl -X POST \
  -H "Content-Type: application/json" \
  -d @inference/example.json
  $INFERENCE_URL/inference
```

## Local testing

To test the example locally you can use [Docker Compose](https://docs.docker.com/compose/install/).

```
# Build the containers locally
docker compose build

# Run the job to set up the data
docker compose up data

# Run the job to train and upload the model
docker compose up training

# Run the inference server
docker compose up inference
```

Access the inference API locally:

```
curl -X POST \
  -H "Content-Type: application/json" \
  -d @inference/example.json
  http://localhost:8080/inference
```
