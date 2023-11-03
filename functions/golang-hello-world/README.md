# Golang hello world

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-go](https://github.com/scaleway/serverless-functions-go) library for local testing.

## Context

This example shows a simple "hello world" function written in Go.

## Setup

Once your environment is set up, you can test your function locally with:

```sh
npm install

go run test/main.go
```

This will launch a local server, allowing you to test the function. For that, you can run in another terminal:

```sh
curl -i -X GET http://localhost:8080
```

If the test succeeded, you can deploy your function with:

```console
serverless deploy
```

