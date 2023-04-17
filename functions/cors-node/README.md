# Node CORS example

This function shows how to handle preflight CORS requests that will be sent by a browser when invoking a function.

The example uses totally permissive, open CORS, you may want to modify this to make it more secure.

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-node](https://github.com/scaleway/serverless-functions-node) library for local testing.

## Context

The function of the example allows most CORS requests. For more documentation on how to control access for CORS requests, see [HTTP response headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers)

## Setup

Once your environment is set up, you can run:

```console
npm install

NODE_ENV=test node handler.js
```

This will launch a local server, allowing you to test the function. For that, you can run in another terminal:

The result should be similar to:

```console
HTTP/1.1 200 OK
access-control-allow-origin: *
access-control-allow-headers: *
access-control-allow-methods: *
access-control-expose-headers: *
content-type: text/plain
content-length: 35
Date: Mon, 17 Apr 2023 07:03:06 GMT
Connection: keep-alive
Keep-Alive: timeout=72
server: envoy

This is allowing most CORS requests
```

You can also check the result of your function in a browser. It should be "This is allowing most CORS requests".

## Deploy and run

Finally, if the test succeeded, you can deploy your function with:

```console
serverless deploy
```

Then, from the given URL, you can run:

```console
# Options request
curl -i -X OPTIONS <function URL>

# Get request
curl -i -X GET <function URL>
```

When invoking this function, the output should be similar to the one obtained when testing locally.
