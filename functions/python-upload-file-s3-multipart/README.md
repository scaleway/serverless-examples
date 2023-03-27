# Python function to upload files to S3

This function does the following steps:

* Read a file from an HTTP request form
* Send the file to long-term storage with Glacier for S3

## Requirements

This example uses the [Python API Framework](https://github.com/scaleway/serverless-api-project) to deploy the function.

If needed, create a bucket and provide the following variables in your environment:

```env
export SCW_ACCESS_KEY =
export SCW_SECRET_KEY =
export BUCKET_NAME =
```

## Running

### Running locally

This examples uses [Serverless Functions Python Framework](https://github.com/scaleway/serverless-functions-python) and can be executed locally:

```console
pip install -r requirements-dev.txt

python app.py
```

The upload endpoint allows you to upload files to Glacier via the `file` form-data key:

```console
echo -e 'Hello world!\n My contents will be stored in a bunker!' > myfile.dat
curl -F file=@myfile.dat localhost:8080
```

### Deploying with the API Framework

Deployment can be done with `scw_serverless`:

```console
scw_serverless deploy app.py
```
