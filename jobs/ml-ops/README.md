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

## How to deploy your MLOps workflow in the cloud?
 
### Step A: Create object storage resources and upload your data file to data store

cf. this [readme](./s3/README.md)

### Step B: Create registry namespaces for your docker images

cf. this [readme](./registry/README.md)

### Step C: Run a machine learning job

cf. this [readme](./job/README.md)

### Step D: Deploy an inference API as a serverless container

cf. this [readme](./container/README.md)
