# Golang CORS Example

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Context

This example shows how to handle preflight CORS requests that will be sent by a browser when invoking a function. The example uses totally permissive, open CORS, you may want to modify this to make it more secure.

## Description

The function of the example is to allow most CORS requests. For more documentation on how to control access for CORS requests, see [HTTP response headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers)

This function uses Golang 1.19 runtime.

## Setup

Once your environment is set up, you can run:

```console
npm install

serverless deploy
```

Then, from the given URL, you can run:

```console
# Options request
curl -i -X OPTIONS <function URL>

# Get request
curl -i -X GET <function URL>
```

The result should be similar to:

```console
HTTP/2 200
access-control-allow-headers: *
access-control-allow-methods: *
access-control-allow-origin: *
content-length: 44
content-type: text/plain
date: Tue, 17 Jan 2023 15:56:52 GMT
x-envoy-upstream-service-time: 10
server: envoy

This function is allowing most CORS requests
```

You can also check the result of your function in a browser. It should be "This function is allowing most CORS requests".
