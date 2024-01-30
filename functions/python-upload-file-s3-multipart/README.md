# Python function to upload files to S3

This function does the following steps:

* Read a file from an HTTP request form
* Store the file in S3

## Requirements

This example uses the [Python API Framework](https://github.com/scaleway/serverless-api-framework-python) to build and deploy the function.

First you need to:

- Create an API key in the [console](https://console.scaleway.com/iam/api-keys), with at least the `ObjectStorageFullAccess` and `FunctionsFullAccess` permissions, and access to the relevant project for Object Storage access
- Get the access key and secret key for this API key
- Get your project ID
- Create an S3 bucket

You then need to set the following environment variables:

```bash
export SCW_ACCESS_KEY=<your access key>
export SCW_SECRET_KEY=<your secret key>
export SCW_DEFAULT_PROJECT_ID=<your project id>
export BUCKET_NAME=<bucket name>
```

## Deploy on Scaleway

Deployment can be done with `scw_serverless`:

```bash
pip install --user -r requirements.txt

scw-serverless deploy app.py
```

This will then print out your function's URL. You can use this to test your function with:

```bash
# Upload the requirements file
curl -F file=@requirements.txt <your function URL>
```

You should then see the `requirements.txt` file uploaded to your bucket.

_Warning_ when deploying the function, do not create a virtualenv directory in this project root, as this will be included in the deployment zip and make it too large.

## Running it locally

You can test your function locally thanks to the [Serverless Functions Python Framework](https://github.com/scaleway/serverless-functions-python). To do this, you can run:

```bash
pip install --user -r requirements-dev.txt

python app.py
```

This starts the function locally, allowing you to upload files to S3 via the `file` form-data key:

```bash
# Upload the requirements file
curl -F file=@requirements.txt localhost:8080
```

