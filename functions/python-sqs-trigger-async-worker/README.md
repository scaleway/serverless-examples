# Connect a function with an async worker via an SQS queue

This example uses a SQS trigger to create a simple front-worker async arch. The front function
receives HTTP requests from the user and sends a message to the queue and returns. Then the SQS trigger
forwards the event to the async worker function at the rate it is able to consume them.

See the documentation on [how to configure a trigger](https://www.scaleway.com/en/docs/serverless/functions/how-to/add-trigger-to-a-function/)
and [how to configure the message retention](https://www.scaleway.com/en/docs/serverless/functions/reference-content/configure-trigger-inputs/) for the queue.

## Deploy

Use the Terraform configuration to deploy the functions, create a SQS queue and attach it with a trigger.

```shell
terraform init
terraform plan -out plan.tfplan && terraform apply plan.tfplan
terraform output -json
```

The Terraform output contains the endpoint of the `front` function.

## Trigger the function

In order to trigger the `worker` function, send messages to the `front` function using the 
endpoint from the Terraform output.

```shell
for i in `seq 1 5`; do
  curl $FRONT_FUNCTION_ENDPOINT
done
```

Looking at the function logs in [Cockpit](https://console.scaleway.com/cockpit/overview) we see that the `front` function is receiving the requests and sending
messages to the SQS queue.

```shell
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
DEBUG - sent source=user stream=stdout
DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
DEBUG - received request source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
DEBUG - sent source=user stream=stdout
DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
DEBUG - received request source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
DEBUG - sent source=user stream=stdout
DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
DEBUG - received request source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
DEBUG - sent source=user stream=stdout
DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
DEBUG - received request source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "curl/7.81.0" source=user stream=stdout
DEBUG - sent source=user stream=stdout
DEBUG - sending notification to queue https://sqs.mnq.fr-par.scw.cloud/ABCCKBASX73H6GHKCLBMP5TLDTV755LZ5XWTKN5XQL2ACE6SCEOXQK75/python-sqs-trigger-hello-world source=user stream=stdout
DEBUG - received request source=user stream=stdout
DEBUG - Function Triggered: / source=core
```

Then the `worker` function logs show that the SQS events are being received:

```shell
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - worker received event: front received event at 2023-08-20 21:33:17.256670 source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - worker received event: front received event at 2023-08-20 21:33:17.111975 source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:17 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - worker received event: front received event at 2023-08-20 21:33:16.953492 source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - worker received event: front received event at 2023-08-20 21:33:16.769310 source=user stream=stdout
DEBUG - Function Triggered: / source=core
DEBUG - 127.0.0.1 - - [20/Aug/2023:21:33:16 +0000] "POST / HTTP/1.1" 200 5 "-" "ScalewayServerlessOrchestrator/1.0" source=user stream=stdout
DEBUG - worker received event: front received event at 2023-08-20 21:33:16.619008 source=user stream=stdout
DEBUG - Function Triggered: / source=core
```
