import os

import asyncio
import boto3
import nats

AWS_ACCESS_KEY_ID = os.environ["AWS_ACCESS_KEY_ID"]
AWS_SECRET_ACCESS_KEY = os.environ["AWS_SECRET_ACCESS_KEY"]
PUBLIC_QUEUE_URL = os.environ["PUBLIC_QUEUE_URL"]
PRIVATE_QUEUE_URL = os.environ["PRIVATE_QUEUE_URL"]

NATS_ENDPOINT_URL = "nats://nats.mnq.fr-par.scaleway.com:4222"
PUBLIC_SUBJECT = os.environ["PUBLIC_SUBJECT"]
PRIVATE_SUBJECT = os.environ["PRIVATE_SUBJECT"]
NATS_CREDS_FILE = os.environ["NATS_CREDS_FILE"]


def send_sqs():
    params = {
        "endpoint_url": "https://sqs.mnq.fr-par.scaleway.com",
        "aws_access_key_id": AWS_ACCESS_KEY_ID,
        "aws_secret_access_key": AWS_SECRET_ACCESS_KEY,
        "region_name": "fr-par",
    }

    sqs = boto3.resource("sqs", **params)

    for queue_url in (PUBLIC_QUEUE_URL, PRIVATE_QUEUE_URL):
        queue = sqs.Queue(queue_url)
        queue_name = queue.attributes["QueueArn"].split(":")[-1]
        print(f"Sending greetings message to SQS {queue_name}...")
        queue.send_message(MessageBody="Hello World SQS!")
        print("Greetings sent!")


async def send_nats():
    nc = await nats.connect(NATS_ENDPOINT_URL, user_credentials=NATS_CREDS_FILE)

    for subject in (PUBLIC_SUBJECT, PRIVATE_SUBJECT):
        print(f"Sending greetings message to NATS {subject}...")
        await nc.publish(subject, b"Hello World NATS!")
        print("Greetings sent!")

    await nc.close()


def main():
    send_sqs()
    asyncio.run(send_nats())


if __name__ == "__main__":
    main()
