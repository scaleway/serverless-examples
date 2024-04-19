# gRPC HTTP2 Server in Go using CLI

Example that demonstrates the deployment of a gRPC service on Scaleway Serverless Containers.

For this example we will use the CLI to deploy the container but you can use [other methods](https://www.scaleway.com/en/docs/serverless/containers/reference-content/deploy-container/).
You can also use the CLI directly from Scaleway console, it saves you operations of setting up your credentials.

## Container settings

To deploy a Serverless Container that uses the gRPC protocol it's important to enable http2 setting (named `h2c` in the API).

## Workflow

Here are the different steps we are going to proceed:

- Quick set-up of Container Registry to host our gRPC container
- Deploy the Serverless Container
- Run gRPC test command ot ensure everything is ok

## Deployment

### Requirements

- [Scaleway CLI](https://github.com/scaleway/scaleway-cli): be sure to be logged in
- **Docker** To build the image
- If you want to test locally, you will need to install the common gRPC stack (doc [here](https://grpc.io/blog/installation/))
  but this is not required for deployment as everything is built in the Dockerfile

### Building the image

To store the image we need to create a namespace in Container Registery with the following command:

```bash
scw registry namespace create name=hello-grpc
```

As output it will give you some informations, keep the endpoint, on this case it's `rg.fr-par.scw.cloud/hello-grpc` (it can change in your case depending the region).

Now open a terminal to set-up Docker, it will allow you to send the built image in the Container Registry we created at the previous step.

Login command:

```bash
docker login rg.fr-par.scw.cloud/hello-grpc -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
```

At this point you have correctly set up docker to be able to push your image online.

**Builb** the image by opening a terminal on your computer in this directory, and run:

```bash
docker build -t grpc:latest .
```

Now we will tag and push the image:

```bash
docker tag grpc:latest rg.fr-par.scw.cloud/hello-grpc/grpc:latest
docker push rg.fr-par.scw.cloud/hello-grpc/grpc:latest
```

### Deploying the image

Now we need to create a nemespace for the Serverless Container and the create the container:

```bash
scw container namespace create name=grpc-test
```

Save the ID of the container.

Now create and deploy the container with a single command:

```bash
scw container container create namespace-id=<PREVIOUS_NAMESPACE_ID> protocol=h2c name=grpc-test registry-image=rg.fr-par.scw.cloud/hello-grpc/grpc:latest
```

Save the DomainName (endpoint) for testing.

### Testing

Before testing, ensure your container is in status `ready`.

In `client/client.go` file, remplace the constant `containerEndpoint` with the DomainName. Do not forget to keep `:80` port at the end even
if you container port is set to 8080, these are two different settings.

And then execute `go run client/client.go' to check if your container reponds.

## Additional content

- [Basic Go gRPC tutorial](https://grpc.io/docs/languages/go/basics/)
