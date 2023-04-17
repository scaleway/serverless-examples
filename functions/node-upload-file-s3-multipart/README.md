# Node function to read files and upload to S3

This example shows how to upload a file to an S3 bucket.

## Requirements

This example assumes that you are familiar with some products of Scaleway's ecosystem:

* how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
* how Object Storage works. Please refer to scaleway's documentation [here](https://www.scaleway.com/en/docs/storage/object/quickstart/) for more information.

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-node](https://github.com/scaleway/serverless-functions-node) library for local testing.

## Context

This example shows how to upload a file to an S3 bucket using serverless function. It also shows how you can test your function locally before deploying.

This function does the following steps:

* Read a file from HTTP request
* Upload file to S3 bucket
* Test locally

## Setup

Ensure to create a bucket and have the following secrets variables available in your environment:

```env
S3_REGION = # Default: fr-par
SCW_ACCESS_KEY =
SCW_SECRET_KEY =
BUCKET_NAME =
```

Once your environment is set up, you can test your function locally with:

```sh
NODE_ENV=test node handler.js
```

This will launch a local server, allowing you to test the function. For that, you can run in another terminal (replace `README.md` with the file you want to upload):

```sh
curl -X POST http://localhost:8080 -H "content-type: multipart/form-data" -F "data=@README.md"
```

The output should be similar to:

```sh
Successfully uploaded README.md to <bucket name>
```

## Deploy and run

Finally, if the test succeeded, you can deploy your function with:

```console
serverless deploy
```

Then, from the given URL, you can run:

```sh
curl -X POST <function URL> -H "Content-Type: multipart/form-data"  -F "data=@README.md"
```

When invoking this function, the output should be similar to the one obtained when testing locally.
