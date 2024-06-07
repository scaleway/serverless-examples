# Rust hello world

This example demonstrates the deployment of a simple rust http service on Scaleway Serverless Containers.

This can be useful if you come from Serverless Functions and you need to install specific dependencies on your system.

For this example, we will use the CLI to deploy the container, but you can use [other methods](https://www.scaleway.com/en/docs/serverless/containers/reference-content/deploy-container/).
You can also use the CLI directly from Scaleway console without having to use your credentials.

## Workflow

Here are the different steps we are going to proceed:

- Quick set-up of Container Registry to host our rust container
- Deploy the Serverless Container
- Test the container

## Deployment

### Requirements

To complete the actions presented below, you must have:
- installed and configured the [Scaleway CLI](https://www.scaleway.com/en/docs/developer-tools/scaleway-cli/quickstart/)
- installed [Docker](https://docs.docker.com/engine/install/) to build the image

### Building the image

1. Run the following command in a terminal to create Container Registry namespace to store the image:

    ```bash
    scw registry namespace create name=hello-rust
    ```

  The registry namespace information displays.

1. Copy the namespace endpoint (in this case, `rg.fr-par.scw.cloud/hello-rust`).

1. Log into the Container Registry namespace you created using Docker:

    ```bash
    docker login rg.fr-par.scw.cloud/hello-rust -u nologin --password-stdin <<< "$SCW_SECRET_KEY"
    ```

  At this point, you have correctly set up Docker to be able to push your image online.

1. In a terminal, access this directory (containing the Dockerfile), and run the following command to build the image:

    ```bash
    docker build -t crabnet:latest .
    ```

1. Tag and push the image to the registry namespace:

    ```bash
    docker tag crabnet:latest rg.fr-par.scw.cloud/hello-rust/crabnet:latest
    docker push rg.fr-par.scw.cloud/hello-ryst/crabnet:latest
    ```

### Deploying the image

In a terminal, run the following command to create a Serverless Containers namespace:

    ```bash
    scw container namespace create name=crabnet
    ```
    The namespace information displays.

1. Copy the namespace ID.

1. Run the following command to create and deploy the container:

    ```bash
    scw container container create namespace-id=<PREVIOUS_NAMESPACE_ID> name=crabnet registry-image=rg.fr-par.scw.cloud/hello-rust/crabnet:latest
    ```
    The container information displays.

1. Copy the DomainName (endpoint) to test your container, you can put the endpoint in your web browser for testing.
