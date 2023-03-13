# Serverless inference with rust functions

Serverless functions are a great choice to deploy inference models! In this example, we deploy a small neural network to recognize hand-written digits.

In addition, this example shows how to load objects from Scaleway's S3 in rust.

You can test the sample application at: <https://rust-example-www.s3-website.fr-par.scw.cloud/>

## Setup

- Install and configure [Cargo](https://doc.rust-lang.org/stable/cargo/getting-started/installation.html)
- [Create an S3 bucket](https://www.scaleway.com/en/docs/storage/object/quickstart/#how-to-create-a-bucket)
- Install and configure [Scaleway's Serverless Framework provider](https://github.com/scaleway/serverless-scaleway-functions#quick-start)
- Setup your local environment:

| Variable | Description |
| :---:   | :---: |
| `SCW_ACCESS_KEY` | Access key to use for S3 operations. |
| `SCW_SECRET_KEY` | Secret key to use for S3 operations. |
| `S3_BUCKET` | Name of the bucket to store the model into. |
| `SCW_DEFAULT_REGION` | Region of your bucket and function. Default: `fr-par`  |

### Training the model

The model is taken from [dfdx examples](https://github.com/coreylowman/dfdx).

If you want to train the model locally you'll need to download the [MNIST dataset](http://yann.lecun.com/exdb/mnist/). The deziped files can then be placed in `./datasets/mnist`.

To train the model and upload it to your bucket, you can run:

```console
cd training && cargo run -r
```

Once finished, the model will be serialized with numpy and uploaded to your S3 bucket.

### Testing locally

A test server is available to make sure your function is working:

```console
cargo run --bin test-server

$ curl localhost:3000 -d '{"data": [0, 1, ...]}'
>> {"output": [0.1233, 0.122, 0.12, ...]}
```

## Running on serverless functions

You can deploy the function directly with Serverless Framework.

```console
serverless deploy
```

## Bonus: running the front-end

If you want to play around with the model, you can try the interface!

Replace `<my-function-url>` with your serverless function endpoint or `localhost:3000` when using the local test server.

```console
export VITE_SLS_FUNCTION_URL=<my-function-url>
pnpm run dev
```
