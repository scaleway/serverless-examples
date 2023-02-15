# Container bash script example

This example shows how to run a Bash script on a schedule using serverless containers.

## Requirements

This example assumes you are familiar with how serverless containers work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/containers/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Example explanation

**Context:** This example shows how to handle scheduled jobs on containers. Each container must have an HTTP server listening on port 8080 to make it possible for the scheduled job to invoke the container.

**Explanation:** In this example, the container launches a bash script every hour. There is a CRON trigger set up in the `serverless.yml` file. In `server.sh` we use `netcat` to listen on port `8080`, and execute a script defined in `script.sh`.

## Setup

Once your environment is set up, you can run:

```console
# Install node dependencies
npm install

# Deploy
serverless deploy
```

## Running

To check things are working, you can either check your container logs to see the job scheduled by the CRON trigger, or you can invoke it directly via an HTTP request (see the deployment output for the container URL). A successful run should show a log containing "Scheduled job executed at $(date).".
