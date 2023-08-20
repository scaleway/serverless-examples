# Connect a front API function with an async worker function using SQS queues and triggers

This example uses a SQS trigger to create a simple front-worker async arch. The front function
receives HTTP requests from the user and sends a message to the queue and returns. Then the SQS trigger
forwards the event to the async worker function at the rate it is able to consume them.

See the documentation to know [how to configure a trigger](https://www.scaleway.com/en/docs/serverless/functions/how-to/add-trigger-to-a-function/)
and the [consideration when configuring the queue message retention](https://www.scaleway.com/en/docs/serverless/functions/reference-content/configure-trigger-inputs/).

## Deploy

Use the Terraform configuration to deploy the functions, create a SQS queue and attach it with a trigger.

```shell
terraform init
terraform plan -out plan.tfplan && terraform apply plan.tfplan
terraform output -json
```

The terraform output contains the endpoint of the `front` function.

## Trigger the function

In order to trigger the `worker` function, send messages to the `front` function using the 
endpoint from the Terraform output.

```shell
for i in `seq 1 5`; do
  curl $FRONT_FUNCTION_ENDPOINT
done
```

Looking at the function logs in Cockpit we see that the `front` function is receiving the requests and sending
messages to the SQS queue.

```shell
2023-08-20 23:33:17.312	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
2023-08-20 23:33:17.312	DEBUG - sent source=user stream=stdout
2023-08-20 23:33:17.312	DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
2023-08-20 23:33:17.312	DEBUG - received request source=user stream=stdout
2023-08-20 23:33:17.266	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:17.174	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
2023-08-20 23:33:17.173	DEBUG - sent source=user stream=stdout
2023-08-20 23:33:17.173	DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
2023-08-20 23:33:17.172	DEBUG - received request source=user stream=stdout
2023-08-20 23:33:17.121	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:17.025	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
2023-08-20 23:33:17.025	DEBUG - sent source=user stream=stdout
2023-08-20 23:33:17.025	DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
2023-08-20 23:33:17.025	DEBUG - received request source=user stream=stdout
2023-08-20 23:33:16.965	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:16.872	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
2023-08-20 23:33:16.872	DEBUG - sent source=user stream=stdout
2023-08-20 23:33:16.872	DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
2023-08-20 23:33:16.871	DEBUG - received request source=user stream=stdout
2023-08-20 23:33:16.780	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:16.685	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
2023-08-20 23:33:16.685	DEBUG - sent source=user stream=stdout
2023-08-20 23:33:16.685	DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
2023-08-20 23:33:16.685	DEBUG - received request source=user stream=stdout
2023-08-20 23:33:16.631	DEBUG - Function Triggered: / source=core
```

Then the `worker` function logs show that the SQS events are being received:

```shell
2023-08-20 23:33:17.319	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-20 23:33:17.319	DEBUG - worker received event: front received event at 2023-08-20 21:33:17.256670 source=user stream=stdout
2023-08-20 23:33:17.316	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:17.181	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-20 23:33:17.180	DEBUG - worker received event: front received event at 2023-08-20 21:33:17.111975 source=user stream=stdout
2023-08-20 23:33:17.177	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:17.040	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-20 23:33:17.040	DEBUG - worker received event: front received event at 2023-08-20 21:33:16.953492 source=user stream=stdout
2023-08-20 23:33:17.037	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:16.859	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-20 23:33:16.858	DEBUG - worker received event: front received event at 2023-08-20 21:33:16.769310 source=user stream=stdout
2023-08-20 23:33:16.855	DEBUG - Function Triggered: / source=core
2023-08-20 23:33:16.718	DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
2023-08-20 23:33:16.717	DEBUG - worker received event: front received event at 2023-08-20 21:33:16.619008 source=user stream=stdout
2023-08-20 23:33:16.713	DEBUG - Function Triggered: / source=core
```
