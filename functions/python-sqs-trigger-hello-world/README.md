# Using SQS queues to trigger functions

It is possible to automatically trigger functions and containers when a SQS queue receives a message.
The message is then received as the body of the usual HTTP request in the `event` parameter.

See the documentation to know [how to configure a trigger](https://www.scaleway.com/en/docs/serverless/functions/how-to/add-trigger-to-a-function/)
and the [consideration when configuring the queue message retention](https://www.scaleway.com/en/docs/serverless/functions/reference-content/configure-trigger-inputs/).

## Deploy

Use the Terraform configuration to deploy the function, create a SQS queue and attach it with a trigger.

```shell
terraform init
terraform plan -out plan.tfplan && terraform apply plan.tfplan
terraform output -json
```

The terraform output contains the values required to send messages to the SQS queue.

## Trigger the function

In order to trigger the function, send messages to the queue using the values from the Terraform output.

```shell
for i in `seq 1 5`; do
AWS_ACCESS_KEY_ID=$ACCESS_KEY AWS_SECRET_ACCESS_KEY=$SECRET_KEY aws sqs send-message \
  --queue-url $URL --endpoint-url $ENDPOINT --region $REGION \
  --message-body "{\"message\": $i}" \
  --output json --no-cli-pager
done
```

Looking at the function logs in Cockpit we see that it is receiving the events sent to SQS.
Since the function sometimes fails, the events are being redelivered until it succeeds, as seen with the
messages 2 and 3.

```shell
2023-08-18 15:24:06.181	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:24:06 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:24:06.181	DEBUG - Success source=user stream=stdout
2023-08-18 15:24:06.181	DEBUG - Received body: {"message": 3} source=user stream=stdout
2023-08-18 15:24:06.178	DEBUG - Function Triggered: / source=core
2023-08-18 15:23:36.166 DEBUG - 127.0.0.1 - - [18/Aug/2023:13:23:36 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:23:36.166	DEBUG - Error source=user stream=stdout
2023-08-18 15:23:36.166	DEBUG - Received body: {"message": 3} source=user stream=stdout
2023-08-18 15:23:36.160	DEBUG - Function Triggered: / source=core
2023-08-18 15:23:06.149	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:23:06 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:23:06.149	DEBUG - Error source=user stream=stdout
2023-08-18 15:23:06.149	DEBUG - Received body: {"message": 3} source=user stream=stdout
2023-08-18 15:23:06.148	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:23:05 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:23:06.147	DEBUG - Success source=user stream=stdout
2023-08-18 15:23:06.147	DEBUG - Received body: {"message": 2} source=user stream=stdout
2023-08-18 15:23:06.063	DEBUG - Function Triggered: / source=core
2023-08-18 15:22:36.152	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:36 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:22:36.152	DEBUG - Success source=user stream=stdout
2023-08-18 15:22:36.152	DEBUG - Received body: {"message": 5} source=user stream=stdout
2023-08-18 15:22:36.147	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:22:36.147	DEBUG - Error source=user stream=stdout
2023-08-18 15:22:36.147	DEBUG - Received body: {"message": 2} source=user stream=stdout
2023-08-18 15:22:36.048	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 19 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:22:36.047	DEBUG - Error source=user stream=stdout
2023-08-18 15:22:36.047 DEBUG - Received body: {"message": 3} source=user stream=stdout
2023-08-18 15:22:35.852	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:22:35.852	DEBUG - Success source=user stream=stdout
2023-08-18 15:22:35.851	DEBUG - Received body: {"message": 1} source=user stream=stdout
2023-08-18 15:22:35.849	DEBUG - 127.0.0.1 - - [18/Aug/2023:13:22:35 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-18 15:22:35.848	DEBUG - Success source=user stream=stdout
2023-08-18 15:22:35.848	DEBUG - Received body: {"message": 4} source=user stream=stdout
2023-08-18 15:22:14.449	DEBUG - [2023-08-18 13:22:14 +0000] [12] [INFO] Booting worker with pid: 12 source=user stream=stderr
2023-08-18 15:22:12.849	DEBUG - Function Triggered: / source=core
```
