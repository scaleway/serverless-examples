# MNQ SQS Publish message using Golang

You can use this example to publish a message to an SQS queue using the MNQ namespace from the Scaleway Go SDK.

## Requirements

This example assumes that you are familiar with some products of Scaleway's ecosystem:

* how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
* how Messaging and Queuing works. Please refer to scaleway's documentation [here](https://www.scaleway.com/en/docs/serverless/messaging/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-go](https://github.com/scaleway/serverless-functions-go) library for local testing.

## Context

This example shows how to create a queue and send a message to a SQS. It can be extended to retrieve or delete queues, send batch messages, delete or retrieve messages, etc...

The function is publishing a message to an SQS queue. The message is a JSON object with the following structure:

```json
{
    "username": "John Doe",
    "message": "Hello World!"
}
```

This function uses Golang 1.18 runtime.

## Setup

To use this example, the following environment variables are needed:

* `SCW_ACCESS_KEY`: Your Scaleway access key
* `SCW_SECRET_KEY`: Your Scaleway secret key
* `SCW_DEFAULT_PROJECT_ID`: Your Scaleway project ID
* `SCW_DEFAULT_REGION`: Your Scaleway region
* `SCW_DEFAULT_ZONE`: Your Scaleway zone
* `SQS_ACCESS_KEY`: Your SQS access key
* `SQS_SECRET_KEY`: Your SQS secret key

Once your environment is set up, you can test your function locally with:

```shell
go run test/main.go
```

This will launch a local server, allowing you to test the function. For that, you can run in another terminal:

```shell
curl -X POST http://localhost:8080 -d '{"username": "John Doe", "message": "Hello World!"}'
```

The status code returned by this request should be `200`.

## Deploy and run

Finally, if the test succeeded, you can deploy your function with:

```console
serverless deploy
```

Then, from the given URL, you can run:

```shell
# POST request
curl -i -X POST <function URL> -d '{"username": "John Doe", "message": "Hello World!"}'
```

Once again, the expected status code for this request is `200`.
