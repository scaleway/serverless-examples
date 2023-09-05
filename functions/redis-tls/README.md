# Redis TLS Example

An example to showcase how to connect a function to a Scaleway Redis cluster with TLS enabled.

At 8:00 AM every day, the function retrieves the hourly temperature in Paris and outputs it to a Redis store.

The Redis certificate is provided via a secret to the function and then written to a file.

## Requirements

This example requires [Terraform](https://www.scaleway.com/en/docs/tutorials/terraform-quickstart/).

## Setup

Everything is managed with Terraform. The terraform config files will also create a Redis cluster, so be sure to remove it from the configuration if you do not need one.

```sh
terraform init
terraform apply
```

You should be able to see your function in the Scaleway console.

## Running

To check the results:

- Get the Redis connection string from the console
- Connect with `redis-cli`:

```sh
# From a Scaleway instance (if acl enabled)
redis-cli -h <cluster_ip> --user <my_user> --askpass --tls --cacert serverless-weather-redis-example.pem
```
