# Title of the function example

**Give a short description/introduction**  
When calling a function from a frontend application in a browser, you may receive an error about missing CORS headers. Here's how you can allow cross-origin requests to a Python function

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Describe the goal of this example

**State the context of the example**  
*This example shows how to handle preflight CORS requests that will be sent by a browser when invoking a function. The example uses totally permissive, open CORS, you will may want to modify this to make it more secure.*

**State what the function does**  
*The function of the example is allowing most CORS requests. For more documentation on how to control access for CORS requests, see [HTTP response headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#the_http_response_headers)*

## Setup

Once your environment is set up, you can run:

```console
npm install

serverless deploy
```

## Running

**Update with the tests to run for your example**  
Then, from the given URL, you can run:

```console
# Options request
curl -i -X OPTIONS <function URL>

# Get request
curl -i -X GET <function URL>
```

**Update with the expected output**  
The result should be similar to:

```console
HTTP/2 200 
access-control-allow-headers: *
access-control-allow-methods: *
access-control-allow-origin: *
content-length: 21
content-type: text/plain
date: Tue, 17 Jan 2023 14:02:46 GMT
server: envoy
x-envoy-upstream-service-time: 222

This is checking CORS%  
```

You can also check the result of your function in a browser. It should be "This is checking CORS".
