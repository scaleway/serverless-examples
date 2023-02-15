# Image labelling example with TensorFlow.js

## Requirements

This example assumes that you are familiar with: 
 * how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
 * how S3 object storage works, and especially how to create a bucket and upload files within a bucket. Please refer to scaleway's documentation [here](https://www.scaleway.com/en/docs/storage/object/quickstart/).
 * how to create and get user API access and secret keys using IAM. Please refer to IAM documentation [here](https://www.scaleway.com/en/docs/identity-and-access-management/iam/concepts/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.


## Context

This example shows how to label an RGB image in an S3 bucket using serverless functions. The example uses a pre-trained ready-to-use model from [TensorFlow.js](https://www.tensorflow.org/js/models). The model is called `mobilenet` and can be used on server or client side. It returns three labels of an image with their respective prediction probabilities (namely, logits). Check `mobilenet` Github repository [here](https://github.com/tensorflow/tfjs-models/tree/master/mobilenet).


## Description

The function gets an RGB image (jpg, jpeg, or png formats) from an S3 bucket, transforms it into a TensorFlow-compatible object (namely, a tensor), and then applys a pretrained `mobilenet` model to return the labels and their respective logits in a json format. 

This function uses Node 18 runtime. Used package are specified in `package.json`.

## Setup

### Upload an image on an S3 bucket 

Create an S3 bucket and upload an RGB image (jpg, jpeg or png format) within this bucket. Keep the name of your bucket and the name of your file (namely, the source key). 

### Fill environment variables

Fill your environment variables within `serverless.yml` file:

```
secret:
    USER_ACCESS_KEY: <bucket scw access key>
    USER_SECRET_KEY: <bucket scw access key id>
    S3_ENDPOINT_URL: http://s3.fr-par.scw.cloud
```

### Install npm modules

Once your environment is set up, you can install `npm` dependencies from `package.json` file using:

```
npm install
```

### Deploy and call the function

Then deploy your function and get its URL using:

```
serverless deploy
```

From the given function URL, you can get image labels in a `json` format:

```
curl -X GET "<function URL>/?sourceBucket=<source bucket name>&sourceKey=<file name within bucket>" | jq
```

The result should be similar to:

```
{
  "labels": [
    {
      "className": "label_1",
      "probability": 0.7794275879859924
    },
    {
      "className": "label_2",
      "probability": 0.11379589140415192
    },
    {
      "className": "label_3",
      "probability": 0.08201524615287781
    }
  ]
}

```

You can also check the result of your function in a browser. It should return the same result.