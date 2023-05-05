# Using Python Requirements with Serverless Framework

If you need to include a PyPI package with your function, you will need to vendor it when packaging your function. This can be done by using `pip install --target package ...` before deploying.

Here's a simple example to achieve that with Serverless Framework.

## Testing with serverless offline for Python

In order to test your function locally before deployment in a serverless function, you can install our offline testing library with:
```
pip install -r requirements-dev.txt
```

Import your environment variables using:

```
export database_model="english-corpus.sqlite3"
```

Launch your function locally:

```
python python-dependencies.py
```

Test your local function using `curl`:
```
curl -i -X POST localhost:8080 -d '{"message":"Hello"}'
```

## Setup

First, you need to set up the [Serverless Framework](https://www.serverless.com/framework/docs/getting-started).

Then, you can run:

```sh
# Install node dependencies
npm install

# Deploy
./bin/deploy.sh
```

The deploy command should print the URL of the new function. You can then use `curl` to check that it works, and should see the response:

```raw
Response status: 200
```
