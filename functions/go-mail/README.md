# Go send mail Example

<!-- To update after tutorial release, can't release both page so this will remain in comment for now: This page is related to this tutorial https://www.scaleway.com/en/docs/tutorials/?facetFilters=%5B%22categories%3Atransactional-email%22%5D&page=1 -->

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Context

This example shows how to send a mail using a Serverless Function via [Scaleway Go SDK](https://github.com/scaleway/scaleway-sdk-go).

## Description

Using Serverless Compute Functions to send mails is a great way to optimise costs, manage automatic scaling and having
advanced logging and metrics per Function.

## Setup

Once your environment is set up, you can test your function locally with:

```sh
npm install
```

## Deploy and run

You can deploy your function with:

```console
serverless deploy
```

Then, from the given URL, you can run:

```console
curl -v -X POST https://YOUR_FUNCTION -H "X-Auth-Token: $SCW_AUTH_TOKEN" --data '{"to": "MAILTO", "subject": "from console test", "message": "very very long msg to api"}'
```

After that you can check your logs in the Scaleway Console.
