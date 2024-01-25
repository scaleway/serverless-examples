# Serverless large messages architecture

This repository contains the source code for the this tutorial: [Create a serverless architecture that manage large messages, with Scaleway Messaging and Queuing NATS, Serverless Functions and Object Storage.](https://github.com/scaleway/docs-content/blob/int-add-mnq-tuto/tutorials/large-messages/index.mdx)

## Requirements

This example assumes that you are familiar with:

- how messaging and queuing works. You can find out more in the [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/messaging/quickstart/)
- how serverless functions work. If needed, you can check [this page](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
- how S3 object storage works, and especially how to create a bucket and upload files within a bucket. Please refer to Scaleway's documentation [here](https://www.scaleway.com/en/docs/storage/object/quickstart/).

## Context

This example shows how to handle large messages using NATS, Object Storage and Serverless Functions. 

## Setup

Once you have cloned this repository you can run:

```console
terraform init
terraform plan
terraform apply
```

## Running

To run the example, you can execute:
```console
./upload_img.sh path/to/your/image
```
Once the script has finished, you should see your image and a pdf with the same name in the bucket you've created.

