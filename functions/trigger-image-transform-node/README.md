# Image transformation with NodeJS

In the tutorial, we transform images stored in an S3 bucket using serverless functions written in Node JS. And SQS trigger to ensure communication between functions

## Requirements

This example assumes that you are familiar with:

* how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
* how S3 object storage works, and especially how to create a bucket and upload files within a bucket. Please refer to scaleway's documentation [here](https://www.scaleway.com/en/docs/storage/object/quickstart/).
* how to create and get user API access and secret keys using IAM. Please refer to IAM documentation [here](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/).
* how Messaging and Queuing works, especially generating credentials [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/messaging/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-node](https://github.com/scaleway/serverless-functions-node) library for local testing.

## Context

This example contains two functions:

  1. Connect to the storage bucket, pull all image files from it, then call the second function to resize each image
  2. Get a specific image (whose name is passed through the call's input data), resize it and push the resized image to a new bucket

## Setup

### Create two buckets

Create an S3 bucket and upload an RGB image (jpg, jpeg or png format) within this bucket. Keep the name of your bucket, this will be your source bucket. Create a second S3 bucket, this will be your destination bucket.

### Create a SQS Message Queue

Ceate a Messaging and Queuing SNS/SQS Namespace and generate credentials using the console, Scaleway CLI or Terraform. 
Using the console or AWS CLI create an SQS queue and keep its name, it will be the queue to which image to tranform will be sent


### Fill environment variables

Ensure to create the buckets, the queue and have the following secrets and environment variables available in your environment (to be able to test locally) and within `serverless.yml` file (to be able to deploy):

```yml
secret:
    S3_ACCESS_KEY_ID: <IAM access key ID with rights over the selected bucket>
    S3_ACCESS_KEY: <IAM access key with rights over the selected bucket>
    SOURCE_BUCKET: <Source bucket>
    DESTINATION_BUCKET: <Destination bucket>
    S3_REGION: <region>
    SQS_ACCESS_KEY: <Generated Messaging and Queueing Access Key>
    SQS_ACCESS_KEY_ID: <Generated Messaging and Queuing Access Key ID>
    QUEUE_URL: <QUEUE_URL>
    SQS_ENDPOINT: <SQS endpoint from Messaging and Queuing
env:
  RESIZED_WIDTH: '500'
```

### Install npm modules

Once your environment is set up, you can install `npm` dependencies from `package.json` file using:

```sh
npm install
```

### Test locally

Once your environment is set up, you can test your functions locally with:

```sh
NODE_ENV=test node BucketScan.js
NODE_ENV=test node ImageTransform.js
```

This will launch local servers, allowing you to test the function. For that, you can run in another terminal:

```sh
curl -X GET "http://localhost:8080"
```

In your logs, you should see something similar to:

```text
Successfully resized <SOURCE_BUCKET>/<filename> and uploaded to <DEST_BUCKET>/resized-<filename>"
```

### Deploy and call the function

Finally, if the test succeeded, you can deploy your function with:

```sh
serverless deploy
```

Once deploy, you will need to add a `SQS Trigger` to the `ImageTransform` function in Scaleway Console
- Select the Namespace where your queue is located
- Type your queue name
- Add the Trigger

Then, from the given BucketScan URL, you can run:
```sh
curl -X GET "<BucketScan function url>"
```

When invoking this function, the output should be similar to the one obtained when testing locally.

Be careful as the `sharp` is dependant on the platform your deploy to you may have to reinstall it before deploying to Scaleway Serverless Functions using

```sh
npm uninstall sharp

npm install sharp --platform=linuxmusl --arch=x64 sharp --ignore-script=false
```
