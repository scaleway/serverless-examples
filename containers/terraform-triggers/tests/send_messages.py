import os

import boto3
import nats

AWS_ACCESS_KEY_ID = os.environ["AWS_ACCESS_KEY_ID"]
AWS_SECRET_ACCESS_KEY = os.environ["AWS_SECRET_ACCESS_KEY"]
PUBLIC_QUEUE_URL = os.environ["PUBLIC_QUEUE_URL"]
PRIVATE_QUEUE_URL = os.environ["PRIVATE_QUEUE_URL"]

NATS_ENDPOINT_URL = "https://nats.mnq.fr-par.scaleway.com"
PUBLIC_SUBJECT = os.environ["PUBLIC_SUBJECT"]
PRIVATE_SUBJECT = os.environ["PRIVATE_SUBJECT"]
NATS_CREDS_FILE = os.environ["NATS_CREDS_FILE"]

params = {
    "endpoint_url": "https://sqs.mnq.fr-par.scaleway.com",
    "aws_access_key_id": AWS_ACCESS_KEY_ID,
    "aws_secret_access_key": AWS_SECRET_ACCESS_KEY,
    "region_name": "fr-par",
}

client = boto3.client("sqs", **params)
sqs = boto3.resource("sqs", **params)

nc = nats.connect(NATS_ENDPOINT_URL, user_credentials=NATS_CREDS_FILE)


def main():
    for queue_url in (PUBLIC_QUEUE_URL, PRIVATE_QUEUE_URL):
        queue = sqs.Queue(queue_url)
        queue_name = queue.attributes["QueueArn"].split(":")[-1]
        print(f"Sending greetings message to SQS {queue_name}...")
        queue.send_message(MessageBody="Hello World SQS!")
        print("Greetings sent!")

    for subject in (PUBLIC_SUBJECT, PRIVATE_SUBJECT):
        nc.publish(subject, b"Hello World NATS!")


if __name__ == "__main__":
    main()
