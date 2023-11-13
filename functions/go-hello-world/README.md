# Golang hello world

This is a simple example of how to create a Golang function to run on Scaleway Serverless Functions.

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-go](https://github.com/scaleway/serverless-functions-go) library for local testing.

## Setup

First set up Serverless Framework with:

```sh
npm i
```

You can run your function locally with:

```sh
go run test/main.go
```

Then in another terminal, you can make a request:

```sh
curl http://localhost:8080
```

... and run the tests:

```sh
go test ./...
```

If the tests succeed, you can deploy your function with:

```sh
serverless deploy
```

Once deployed, you can again submit a request using `curl` to the URL printed in the deployment output.
