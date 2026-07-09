# pg-connection-retry

A minimal Node.js example demonstrating how to make a Scaleway Serverless Container resilient to VPC connection delays when it starts.

## Context

Sometimes, when a Serverless Container is reloaded, it spawns **before** its VPC connection is fully ready. This is a known issue on the VPC side. Until it is fixed, applications need to **retry** their database connection on startup rather than crashing immediately.

This example shows a simple Express app that:

1. Tries to connect to a Postgres database on startup.
2. Retries every 5 seconds (up to 12 times) if the connection fails.
3. Only starts the HTTP server once a connection is established.

## How it works

- The app uses the [`pg`](https://node-postgres.com/) library and a standard libpq connection (via `PGHOST`, `PGPORT`, `PGDATABASE`, `PGUSER`, `PGPASSWORD` environment variables — no hardcoded credentials).
- `waitForDatabase()` runs a `SELECT 1` query in a retry loop before starting the Express server.
- If all retries are exhausted, the process exits with a non-zero code so the container platform can restart it.

## Files

| File           | Description                                                                    |
|----------------|--------------------------------------------------------------------------------|
| `index.js`     | The Node.js application with retry logic                                       |
| `package.json` | Dependencies (`express`, `pg`)                                                 |
| `Dockerfile`   | Builds the container image                                                     |
| `main.tf`      | Terraform config: RDB instance, VPC, private network, and Serverless Container |
| `variables.tf` | Terraform input variables                                                      |
| `providers.tf` | Terraform provider configuration                                               |
| `outputs.tf`   | Terraform outputs (container endpoint)                                         |

## Prerequisites

- A Scaleway account
- [Terraform](https://developer.hashicorp.com/terraform/install) or [OpenTofu](https://opentofu.org/) installed
- [Docker](https://docs.docker.com/get-docker/) installed
- [Scaleway CLI](https://www.scaleway.com/en/docs/scaleway-cli/) installed and authenticated

## Deploy

### 1. Build and push the container image

Log in to the Scaleway Container Registry:

```bash
scw registry login
```

Build and push the image:

```bash
docker build -t rg.fr-par.scw.cloud/examples/pg-connection-retry:latest .
docker push rg.fr-par.scw.cloud/examples/pg-connection-retry:latest
```

### 2. Deploy with Terraform

Set the required variables (e.g. via environment variables):

```bash
export TF_VAR_db_admin_password="your-secure-admin-password"
export TF_VAR_db_password="your-secure-app-password"
export TF_VAR_registry_image="rg.fr-par.scw.cloud/examples/pg-connection-retry:latest"
```

Apply the Terraform configuration:

```bash
tofu init
tofu apply
```

### 3. Access the container

After the deployment, you can find the container endpoint in the output:

```bash
tofu output container_public_endpoint
```

You can test the connection:

```bash
CONTAINER_URL=$(tofu output -raw container_public_endpoint)

# Health check
curl $CONTAINER_URL/

# Database ping
curl $CONTAINER_URL/ping
```

## Cleaning Up

Destroy the Terraform-managed infrastructure when you're done:

```bash
tofu destroy
```
