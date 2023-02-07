# MNQ SQS Publish Example using Golang

You can use this example to publish a message to an SQS queue using the MNQ namespace from the Scaleway Go SDK.

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Description

The function of the example is publishing a message to an SQS queue. The message is a JSON object with the following structure:

```json
{
    "username": "John Doe",
    "message": "Hello World!"
}
```

This function uses Golang 1.18 runtime.

## Setup

Once your environment is set up, you can run:

```console
go run .

serverless deploy
```

## Needed environment variables

- `SCW_ACCESS_KEY`: Your Scaleway access key
- `SCW_SECRET_KEY`: Your Scaleway secret key
- `SCW_DEFAULT_PROJECT_ID`: Your Scaleway project ID
- `SCW_DEFAULT_REGION`: Your Scaleway region
- `SCW_DEFAULT_ZONE`: Your Scaleway zone

- `SQS_ACCESS_KEY`: Your SQS access key
- `SQS_SECRET_KEY`: Your SQS secret key

## Running

Then, from the given URL, you can run:

```shell
# Options request
curl -i -X POST <function URL> -d '{"username": "John Doe", "message": "Hello World!"}'
```
