# Node CORS example

This function shows how to handle preflight CORS requests that will be sent by a browser when invoking a function.

The example uses totally permissive, open CORS, you may want to modify this to make it more secure.

## Setup

This example uses the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions).

Once this is set up, you can run:

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
