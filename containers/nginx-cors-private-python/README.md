# Using NGINX as a proxy for a private container to handle CORS

Demonstration using [Scaleway Serverless Containers](https://www.scaleway.com/en/serverless-containers/) from the browser to handle normal API requests and serve images.

## Setup

- Install and configure [Terraform](https://developer.hashicorp.com/terraform/tutorials/certification-associate-tutorials/install-cli)
- Install the [Scaleway CLI](https://github.com/scaleway/scaleway-cli#installation)
- Install [Docker compose](https://docs.docker.com/compose/) to test locally

## Running on serverless containers

Log in to registry, create a namespace and build and push the images:

```
make docker-login
make create-namespace

make build-gateway push-gateway
make build-server push-server
```

Set up Terraform with:

```
make tf-init
```

Run plan, have a look at what it will created:

```
make tf-plan
```

Deploy with:

```
make tf-apply
```

This will template a file at `index.html` in the root of the project, with the function URL and token populated. You can then open `index.html` in your browser, e.g.

```
firefox index.html
```

## Local test

Test the setup locally by running the deploy, then modifying `functionUrl` in `index.html` to `localhost:8080`.

You can then run:

```
docker-compose up
```

Then open `index.html` in your browser, e.g.

```
firefox index.html
```

