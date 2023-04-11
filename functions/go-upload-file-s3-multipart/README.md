# Go function to read file and upload it to S3

This example shows how to upload a file to an S3 bucket.

## Requirements

This example assumes that you are familiar with some products of Scaleway's ecosystem:

* how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
* how Object Storage works. Please refer to scaleway's documentation [here](https://www.scaleway.com/en/docs/storage/object/quickstart/) for more information.

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-go](https://github.com/scaleway/serverless-functions-go) library for local testing.

## Context

This example shows how to upload a file to an S3 bucket using serverless function. It also shows how you can test your function locally before deploying.

This function does the following steps:

* Read a file from HTTP request
* Save the file locally in ephemeral storage
* Send file to S3 bucket

This function uses Golang 1.20 runtime.

## Setup

If you want to enable S3 upload, ensure to create a bucket and have the following secrets variables available in your environment:

```sh
S3_ENABLED=true
S3_ENDPOINT= # ex: sample.s3.fr-par.scw.cloud
S3_ACCESSKEY=
S3_SECRET=
S3_BUCKET_NAME=
S3_REGION= # ex: fr-par
```

If s3 is not enabled, the file will be saved on the ephemeral storage of your function.

Once your environment is set up, you can test your function locally with:

```sh
go run test/main.go
```

This will launch a local server, allowing you to test the function. For that, you can run in another terminal (replace `go.sum` with the file you want to upload):

```sh
curl -X POST http://localhost:8080 -H "Content-Type: multipart/form-data" -F "data=@go.sum"
```

The logs should be similar to:

```sh
2023/04/11 11:30:06 S3 upload enabled
2023/04/11 11:30:08 successfully created upload-file-to-s3
2023/04/11 11:30:09 Successfully uploaded /var/folders/wn/qnp2ebt54mz040bgffg35xt80000gn/T/go.sum of size 1234
```

## Deploy and run

Finally, if the test succeeded, you can deploy your function with:

```console
serverless deploy
```

Then, from the given URL, you can run:

```sh
curl -X POST <function URL> -H "Content-Type: multipart/form-data"  -F "data=@go.sum"
```

When invoking this function, the output should be similar to the one obtained when testing locally.
