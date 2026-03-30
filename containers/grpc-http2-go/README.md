# gRPC HTTP2 Server in Go using CLI

This example demonstrates the deployment of a gRPC service on Scaleway Serverless Containers.

For this example, we will use the CLI to deploy the container, but you can use [other methods](https://www.scaleway.com/en/docs/serverless/containers/reference-content/deploy-container/).
You can also use the CLI directly from Scaleway console without having to use your credentials.

## Workflow

Here are the different steps we are going to proceed:

- Quick set-up of Container Registry to host our gRPC container
- Deploy the Serverless Container
- Run gRPC test command ot ensure everything is ok

## Deployment

### Requirements

To complete the actions presented below, you must have:

- installed and configured the [Scaleway CLI](https://www.scaleway.com/en/docs/developer-tools/scaleway-cli/quickstart/)
- installed [Docker](https://docs.docker.com/engine/install/) to build the image
- installed the common [gRPC stack](https://grpc.io/blog/installation/)) to test locally (optional)

### Building the image

- Run the following command in a terminal to create Container Registry namespace to store the image:

```bash
scw registry namespace create name=hello-grpc
```
  
- Copy the namespace endpoint (in this case, `rg.fr-par.scw.cloud/hello-grpc`).

- Login to the Scaleway Container Registry using the CLI:

```bash
scw registry login
```

At this point, you have correctly set up Docker to be able to push your image online.

In a terminal, access this directory (containing the Dockerfile), and run the following command to build the image:

```bash
docker build -t grpc:latest .
```

- Tag and push the image to the registry namespace:

```bash
docker tag grpc:latest rg.fr-par.scw.cloud/hello-grpc/grpc:latest
docker push rg.fr-par.scw.cloud/hello-grpc/grpc:latest
```

### Deploying the image

In a terminal, run the following command to create a Serverless Containers namespace:

```bash
scw container namespace create name=grpc-test
```

The namespace information displays.

- Copy the namespace ID.
- Run the following command to create and deploy the container (make sure to use the `h2c` protocol to connect via HTTP2):

```bash
scw container container create namespace-id=<PREVIOUS_NAMESPACE_ID> protocol=h2c name=grpc-test registry-image=rg.fr-par.scw.cloud/hello-grpc/grpc:latest
```

- Copy the DomainName (endpoint) to test your container.

### Testing

Make sure your container is in a `ready` status before testing it.

- In the `client/client.go` file, replace the constant `containerEndpoint` with the `DomainName` copied previously.

- Make sure to keep the `:443` port at the end even if you container port is set to 8080, as these are two different settings.

- Run the command below to check if your container responds:

```bash
go run client/client.go
```

## Additional content

- [Basic Go gRPC tutorial](https://grpc.io/docs/languages/go/basics/)
