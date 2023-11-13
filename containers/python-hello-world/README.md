# Python Hello World

This demonstrates a simple example of running a [Flask](https://flask.palletsprojects.com/en/2.0.x/quickstart/) HTTP server in a [Scaleway Serverless Container](https://www.scaleway.com/en/serverless-containers/).

## Requirements

This example assumes you are familiar with how Serverless Containers work. If needed, you can
check the [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/containers/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Deployment

Once your environment is set up, you can deploy your container with:

```sh
npm i

serverless deploy
```

When the deployment is complete, you should be able to `curl` the container's endpoint or hit it from a browser and see the message returned by the Flask server.
