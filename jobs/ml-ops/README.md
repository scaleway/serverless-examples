# Serverless MLOps

In this example, we train and deploy a binary classification inference API using serverless computing resources (job+container). We use object storage resources to store data and training artifacts. We use container registry to store docker images.

## Use case: Bank Telemarketing

### Context

We use a bank telemarketing dataset to predict if a client would engage in a term deposit subscription. This dataset records marketing phone calls made to clients. The outcome of the phone call is in shown in the `y` column:
* `0` : no subscription
* `1` : subscription

### Data Source

The dataset has many versions and is open-sourced and published [here](http://archive.ics.uci.edu/dataset/222/bank+marketing) on the UCI Machine Leaning repository and is close to the one analyzed in the following research work:

* [Moro et al., 2014] S. Moro, P. Cortez and P. Rita. A Data-Driven Approach to Predict the Success of Bank Telemarketing. Decision Support Systems, Elsevier, 62:22-31, June 2014

We use the dataset labelled in the source as `bank-additional-full.csv`. You can download, extract this file, rename it to `bank_telemarketing.csv` then put it under this [directory](./s3/data-store/data/).

## How to deploy your MLOps pipeline on Scaleway Cloud?
 
### Step A: Create cloud resources for the ML pipeline

Create `.env` file in `jobs/data-loader-job` and `jobs/ml-job` directories and fill them as it follows:

```text
SCW_ACCESS_KEY=<access-key>
SCW_SECRET_KEY=<secret-key>
```

Create `.tfvars` file in `/terraform` directory and put variable values in it:

```
region = "fr-par"
access_key = "<access-key>"
secret_key = "<secret_key>"
project_id = "<project_id>"
data_file     = "bank_telemarketing.csv"
model_object  = "classifier.pkl"
image_version = "v1"
```

Then perform:

```bash
cd terraform
terraform init
terraform plan -var-file=testing.tfvars
terraform apply -var-file=testing.tfvars
```

### Step B: Define and run a job to ship data from public source to s3

Use the console to define and run the data loader job using image pushed to Scaleway registry.

cf. this [readme](./jobs/data-loader-job/README.md)

### Step C: Define and run the ML job to train classifier

Use the console to define and the ML job using image pushed to Scaleway registry.

cf. this [readme](./jobs/ml-job/README.md)

### Step D: Call your serverless container to (re)load model and to get inference results 

cf. this [readme](./containers/inference-api/README.md)
