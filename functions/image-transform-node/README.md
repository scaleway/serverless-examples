With this tutorial you will learn how to transform images in a S3 bucket using FaaS and Node JS.
Serverless Function can help you schedule automated transformation of all jpg or png pictures stored in an Object Storage Bucket.

Detailed Tutorial can be found on Scaleway's [Documentation Website](https://www.scaleway.com/en/docs/tutorials/transform-images-using-faas-and-nodejs/?facetFilters=%5B%22categories%3Afunctions%22%5D&page=1)


# Project Description

This project contains two functions where:
  * The first connects to the source Object Storage bucket, gets all image files from it, then calls the second function to resize each image.
  * The second function gets a specific image (whose name is passed through an HTTP call), resizes it and pushes the resized image to a new bucket.

# Requirements

  - You have an account and are logged into the [Scaleway console](https://console.scaleway.com)
  - You have created a [Serverless Function namespace](https://www.scaleway.com/en/docs/compute/functions/how-to/create-a-functions-namespace/)
  - You have created two [Object Storage buckets](https://www.scaleway.com/en/docs/storage/object/how-to/create-a-bucket/) (one which contains the source pictures and the other for the destination pictures)
  - You have generated your [Scaleway API keys](https://www.scaleway.com/en/docs/console/my-project/how-to/generate-api-key/)
  - **(Optional)** You have installed the [Serverless Framework](https://serverless.com/) with Scaleway's Serverless Plugin


# Run the project
    - Run`npm i`to install required dependencies 
    - Be careful to install sharp using _linuxmusl_ platform to ensure compatibility with scaleway's nodejs runtime
        `npm install --platform=linuxmusl --arch=x64 sharp --ignore-script=false --save `

