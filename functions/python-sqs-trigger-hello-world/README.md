# Using SQS queues to trigger functions

It is possible to automatically trigger functions and containers when a SQS queue receives a message.
The message is then received as the body of the usual HTTP request in the `event` parameter.

See the documentation on [how to configure a trigger](https://www.scaleway.com/en/docs/serverless/functions/how-to/add-trigger-to-a-function/)
and [how to configure the message retention](https://www.scaleway.com/en/docs/serverless/functions/reference-content/configure-trigger-inputs/) for the queue.

## Deploy

Use the Terraform configuration to deploy the function, create a SQS queue and attach it with a trigger.

```shell
terraform init
terraform plan -out plan.tfplan && terraform apply plan.tfplan
terraform output -json
```

The Terraform output contains the values required to send messages to the SQS queue.

## Trigger the function

To trigger the function, send messages to the queue using the values from the Terraform output.

```shell
export AWS_ACCESS_KEY_ID=$ACCESS_KEY
export AWS_SECRET_ACCESS_KEY=$SECRET_KEY

for i in `seq 1 5`; do
aws sqs send-message \
  --queue-url $URL \
  --endpoint-url $ENDPOINT \
  --region $REGION \
  --message-body "{\"message\": $i}" \
  --output json \
  --no-cli-pager
done
```

Looking at the function logs in [Cockpit](https://console.scaleway.com/cockpit/overview) we see that it is receiving the events sent to SQS.
Since the function sometimes fails, the events are being redelivered until it succeeds, as seen with the
messages 2 and 3.

```shell
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:24:06 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Success source=user stream=stdout
DEBUG - Received body: {"message": 3} source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:23:36 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Error source=user stream=stdout
DEBUG - Received body: {"message": 3} source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:23:06 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Error source=user stream=stdout
DEBUG - Received body: {"message": 3} source=user stream=stdout
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:23:05 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Success source=user stream=stdout
DEBUG - Received body: {"message": 2} source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:36 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Success source=user stream=stdout
DEBUG - Received body: {"message": 5} source=user stream=stdout
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Error source=user stream=stdout
DEBUG - Received body: {"message": 2} source=user stream=stdout
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Error source=user stream=stdout
DEBUG - Received body: {"message": 3} source=user stream=stdout
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Success source=user stream=stdout
DEBUG - Received body: {"message": 1} source=user stream=stdout
DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - Success source=user stream=stdout
DEBUG - Received body: {"message": 4} source=user stream=stdout
DEBUG - [2023-08-18 13:22:14 +0000] [12] [INFO] Booting worker with pid: 12 source=user stream=stderr
DEBUG - Function Triggered: / source=core
```
