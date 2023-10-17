# Scaleway Serverless Examples

This is a collection of example projects and patterns for use with Scaleway's serverless products.

Useful Links:

- [Function Documentation](https://www.scaleway.com/en/docs/serverless/functions/)
- [Containers Documentation](https://www.scaleway.com/en/docs/serverless/containers/)
- [Slack Community][slack-scaleway] [#serverless-functions][slack-functions] [#serverless-containers][slack-containers]

[slack-scaleway]: https://slack.scaleway.com/
[slack-functions]: https://scaleway-community.slack.com/app_redirect?channel=serverless-functions
[slack-containers]: https://scaleway-community.slack.com/app_redirect?channel=serverless-containers

Table of Contents:

- [Scaleway Serverless Examples](#scaleway-serverless-examples)
  - [Examples](#examples)
    - [🚀 Functions](#-functions)
    - [📦 Containers](#-containers)
    - [💜 Projects](#-projects)
  - [Contributing](#contributing)

## Examples

### 🚀 Functions

<!-- markdownlint-disable MD033 -->
| Example                                                                                                                                                                               | Runtime   | Deployment             |
|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|------------------------|
| **[Badge PHP](functions/badge-php/README.md)** <br/> A PHP function to generate repository badges.                                                                                    | php82     | [Serverless Framework] |
| **[CORS Go](functions/cors-go/README.md)** <br/> A Go function which allows CORS requests.                                                                                            | go119     | [Serverless Framework] |
| **[CORS Node](functions/cors-node/README.md)** <br/> A Node function which allows CORS requests.                                                                                      | node18    | [Serverless Framework] |
| **[CORS Python](functions/cors-python/README.md)** <br/> A Python function which allows CORS requests.                                                                                | python310 | [Serverless Framework] |
| **[CORS Rust](functions/cors-rust/README.md)** <br/> A Rust function which allows CORS requests.                                                                                      | rust165   | [Serverless Framework] |
| **[Go SQS Publish](functions/go-mnq-sqs-publish/README.md)** <br/> A Go function to publish messages to SQS.                                                                          | go118     | [Serverless Framework] |
| **[Go MultiPart Upload to S3](functions/go-upload-file-s3-multipart)** <br/> A function to upload file from form-data to S3.                                                          | go120     | [Serverless Framework] |
| **[Image Transform](functions/image-transform-node/README.md)** <br/> A function that resizes images from an S3 bucket.                                                               | node16    | [Serverless Framework] |
| **[Image Transform with triggers](functions/trigger-image-transform-node/README.md)** <br/> A function that resizes images from an S3 bucket and use SQS triggers to smooth traffic.  | node20    | [Serverless Framework] |
| **[Node MultiPart Upload to S3](functions/node-upload-file-s3-multipart/README.md)** <br/> A function to upload file from form-data to S3.                                            | node19    | [Serverless Framework] |
| **[PHP write to S3](functions/php-s3/README.md)** <br/> A PHP function that connects to, and writes to an S3 bucket.                                                                  | php82     | [Terraform]            |
| **[PostgeSQL Node](functions/postgre-sql-node/README.md)** <br/> A Node function to connect and interact with PostgreSQL database.                                                    | node18    | [Serverless Framework] |
| **[Python ChatBot](functions/python-dependencies/README.md)** <br/> A chatbot example with ChatterBot.                                                                                | python310 | [Serverless Framework] |
| **[Python Dependencies](functions/python-chatbot/README.md)** <br/> Example showing how to use Python requirements with Serverless Framework.                                         | python310 | [Serverless Framework] |
| **[Python MultiPart Upload to S3](functions/python-upload-file-s3-multipart/README.md)** <br/> A function to upload file from form-data to S3.                                        | python311 | [Python API Framework] |
| **[Python SQS Trigger Hello World](functions/python-sqs-trigger-hello-world/README.md)** <br/> Trigger a function by sending a message to a SQS queue.                                | python311 | [Terraform]            |
| **[Python SQS Trigger Async Worker](functions/python-sqs-trigger-async-worker/README.md)** <br/> Use SQS queues and function triggers to scheule an async task from another function. | python311 | [Terraform]            |
| **[Redis TLS](functions/redis-tls/README.md)** <br/> How to connect a function to a Scaleway Redis cluster with TLS enabled.                                                          | python310 | [Terraform]            |
| **[Rust MNIST](functions/rust-mnist/README.md)** <br/> A Rust function to recognize hand-written digits with a simple neural network.                                                 | rust165   | [Serverless Framework] |
| **[PostgreSQL Python](functions/postgre-sql-python/README.md)** <br/> A Python function to perform a query on a PostgreSQL managed database.                                          | python310 | [Serverless Framework] |
| **[Terraform Python](functions/terraform-python-example/README.md)** <br/> A Python function deployed with Terraform.                                                                 | python310 | [Terraform]            |
| **[Triggers Getting Started](functions/triggers-getting-started/README.md)** <br/> Simple SQS trigger example for all runtimes.                                                       | all       | [Terraform]            |
| **[Typescript with Node runtime](functions/typescript-with-node/README.md)** <br/> A Typescript function using Node runtime.                                                          | node18    | [Serverless Framework] |
| **[Serverless Gateway Python Example](functions/serverless-gateway-python/README.md)** <br/> A Python serverless API using Serverless Gateway.                                        | python310 | [Python API Framework] |

### 📦 Containers

| Example                                                                                                                                                                                     | Language     | Deployment             |
|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------|------------------------|
| **[Container Bash Script](containers/bash-scheduled-job/README.md)** <br/> A Bash script runnning on a schedule using serverless containers.                                                | Bash         | [Serverless Framework] |
| **[Function Handler Java](containers/function-handler-java/README.md)** <br/> A Java function handler deployed on CaaS.                                                                     | Java         | [Serverless Framework] |
| **[NGINX CORS Private](containers/nginx-cors-private-python/README.md)** <br/> An NGINX proxy to allow CORS requests to a private container.                                                | Python Flask | [Terraform]            |
| **[NGINX hello world](containers/nginx-hello-world/README.md)** <br/> A minimal example running the base NGINX image in a serverless container.                                             | N/A          | [Serverless Framework] |
| **[Terraform NGINX hello world](containers/terraform-nginx-hello-world/README.md)** <br/> A minimal example running the base NGINX image in a serverless container deployed with Terraform. | N/A          | [Terraform]            |

### 💜 Projects

| Example                                                                                                                                   | Services    | Language | Deployment             |
|-------------------------------------------------------------------------------------------------------------------------------------------|-------------|----------|------------------------|
| **[Kong API Gateway](projects/kong-api-gateway/README.md)** <br/> Deploying a Kong Gateway on containers to provide routing to functions. | CaaS & FaaS | Python   | [Serverless Framework] |
| **[Serverless Gateway](https://github.com/scaleway/serverless-gateway)** <br/> Our serverless gateway for functions and containers.       | API Gateway | Python   | [Python API Framework] |
| **[Monitoring Glaciers](projects/blogpost-glacier/README.md)** <br/> A project to monitor glaciers and the impact of global warming.      | S3 & RDB    | Golang   | [Serverless Framework] |

[Serverless Framework]: https://github.com/scaleway/serverless-scaleway-functions
[Terraform]: https://registry.terraform.io/providers/scaleway/scaleway/latest/docs
[Python API Framework]: https://github.com/scaleway/serverless-api-project

## Contributing

Want to share an example with the community? 🚀

We'd love to accept your contributions. Here are the steps to provide an example:

- [Fork the repository](https://github.com/scaleway/serverless-examples/fork)
- Write your example.
- Add a README.md. Please use [the provided template](docs/templates/readme-example-template.md).
- Open a new [Pull Request](https://github.com/scaleway/serverless-examples/compare).

You can also [open an issue](https://github.com/scaleway/serverless-examples/issues/new) to suggest new examples or improvements to existing ones!
