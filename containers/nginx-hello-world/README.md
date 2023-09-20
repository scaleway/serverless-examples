# NGINX hello world

This demonstrates a simple example of running the [NGINX base image](https://hub.docker.com/_/nginx/) on [Scaleway Serverless Containers](https://www.scaleway.com/en/serverless-containers/).

## Requirements

This example assumes you are familiar with how serverless containers work. If needed, you can
check the [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/containers/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Deployment

Once your environment is set up, you can test your function locally with:

```sh
npm install

serverless deploy
```
