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

Edit `terraform/vars/main.tfvars` to set your project ID:

```
# Replace with your project ID
project_id = "12ef4x91-yh12-1234-g22g-83er2q4z51ec"
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

This will generate a script at `curl.sh` which will use cURL to invoke the private function directly to test:

```
./curl.sh
```

This will also generate `index.html` in the root of the project, which you can open in your browser, e.g.

```
firefox index.html
```

This will then make CORS requests to the public function, which will respond properly, and forward other requests to the private container.

## Local test

Test the setup locally by running the deploy, then modifying `functionDomain` in `index.html` to `localhost:8080`.

You can then run:

```
docker compose up
```

Then open `index.html` in your browser, e.g.

```
firefox index.html
```
