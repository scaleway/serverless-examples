# Using Python Requirements with Serverless Framework

If you need to include a PyPI package with your function, you will need to include it in the zip file packaged with your function.

This can be done by manually using `pip install --target package ...` to install dependencies in a directory, then zipping this up with your code. You can also use the [`serverless-python-requirements` plugin](https://github.com/serverless/serverless-python-requirements) for Serverless Framework.

## Requirements

This example assumes that you are familiar with how Serverless Functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

> [!NOTE]
> You must also have the same version of Python used in your function available locally on your system, i.e. if you are using the function runtime `python310`, you must have the command `python3.10` available in your shell.

## Running on Scaleway

Run the following:

```console
# Install node dependencies
npm install

# Deploy
serverless deploy
```

The deploy command should print the URL of the new function.

You can then use `curl` to check that it works, and should see the response:

```raw
Response status: 200
```

## Testing with serverless offline for Python

In order to test your function locally before deployment in a serverless function, you can use our offline testing library with:

```console
pip install -r requirements-dev.txt
```

Launch your function locally:

```console
python3 handler.py
```

Test your local function using `curl` (in another terminal)

```console
curl localhost:8080
```

