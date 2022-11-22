# Serverless inference with rust functions

Serverless functions are a great choice to deploy inference models! In this example, we deploy a small neural network to recognize hand-written digits.

In addition, this example shows how to load objects from Scaleway's S3 in rust.

You can test the sample application on: <https://rust-example-www.s3-website.fr-par.scw.cloud/>

## Setup

- Install and configure [Cargo](https://doc.rust-lang.org/stable/cargo/getting-started/installation.html)

### Training the model

The model is taken from [dfdx examples](https://github.com/coreylowman/dfdx).

If you want to train the model locally you'll need to download the [MNIST dataset](http://yann.lecun.com/exdb/mnist/). The deziped files can then be placed in `datasets/mnist`.

Once finished, the model will be serialized with numpy and uploaded to your S3 bucket.

### Testing locally

A test server is available to play around with your function:

```console
cargo run --bin test-server

curl localhost:3000 -d '{"data": [0, 1, ...]}'
```

## Running on serverless functions

```console
serverless deploy
```
