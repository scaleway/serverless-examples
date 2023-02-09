# Container bash script example

This example shows how to handle scheduled jobs on containers.

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Example explanation

**Context:** This example shows how to handle scheduled jobs on containers. Each container must have an HTTP server listening on port 8080 to make it possible for the scheduled job to invoke the container.

**Explanation:** In this example, the container launches a bash script every hour : there is a CRON trigger set up in the serverless.yml file. The bash script (file server.sh) uses netcat to listen on port 8080 and executes a script saved in the file script.sh.

## Setup

Once your environment is set up, you can run:

```console
# Install node dependencies
npm install

# Deploy
serverless deploy
```

## Running

Then, you can check your scheduled job by checking the logs. Every hour, a new log containing "Scheduled job executed at $(date)." should appear.
