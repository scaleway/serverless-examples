# Go function to read file and upload to S3

This function does the following steps :
* Read a file from HTTP request
* Save the file locally in ephemeral storage
* Send file to S3 bucket

Additionnaly it uses the [serverless-functions-go](https://github.com/scaleway/serverless-functions-go) library for local testing

## Requirements

If you want to enable S3 upload, ensure to create a bucket and have the following secrets variables available in your environment:

```
S3_ENABLED=true
S3_ENDPOINT= # ex: sample.s3.fr-par.scw.cloud
S3_ACCESSKEY=
S3_SECRET=
S3_BUCKET_NAME=
S3_REGION= # ex: fr-par
```

If s3 is not enabled the file will be saved on the ephemeral storage of your function.

```sh
go get

# for local testing :
go run cmd/main.go
```

To call the function (replace `go.sum` with the file you want to upload):

```sh
curl -X POST -H "Content-Type: multipart/form-data"  -F "data=@go.sum" http://localhost:8080
```
