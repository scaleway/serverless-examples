# Example showing how to use Python requirements with Serverless Framework

First you need to setup the [Serverless Framework](https://www.serverless.com/framework/docs/getting-started).

Then, you can run:

```
# Install node dependencies
npm install

# Deploy
./bin/deploy.sh
```

The deploy command should print the URL of the new function. You can then use `curl` to check that it works, and should see the response:

```
Response status: 200
```
