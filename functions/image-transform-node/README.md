# Image transformation with NodeJS

This code is from the tutorial featured on Scaleway's [documentation website](https://www.scaleway.com/en/docs/tutorials/transform-images-using-faas-and-nodejs/).

In the tutorial, we transform images stored in an S3 bucket using serverless functions written in Node JS.

## Functions

This example contains two functions:

  1. Connect to the storage bucket, pull all image files from it, then call the second function to resize each image
  2. Get a specific image (whose name is passed through the call's input data), resize it and push the resized image to a new bucket

## Setup and run

You can follow the instructions in the [tutorial](https://www.scaleway.com/en/docs/tutorials/transform-images-using-faas-and-nodejs/) to set up and run the code.
