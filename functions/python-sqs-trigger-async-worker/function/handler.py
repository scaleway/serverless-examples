import datetime
import os

import boto3


def handle_front(event, context):
    print('received request')

    sqs_access_key = os.environ['SQS_ACCESS_KEY']
    sqs_secret_key = os.environ['SQS_SECRET_KEY']
    sqs_endpoint = os.environ['SQS_ENDPOINT']
    sqs_queue_url = os.environ['SQS_QUEUE_URL']
    sqs_region = os.environ['SQS_REGION']

    sqs = boto3.client(
        'sqs',
        region_name=sqs_region,
        use_ssl=True,
        endpoint_url=sqs_endpoint,
        aws_access_key_id=sqs_access_key,
        aws_secret_access_key=sqs_secret_key,
    )

    print(f'sending notification to queue {sqs_queue_url}')
    try:
        sqs.send_message(
            QueueUrl=sqs_queue_url,
            MessageBody=f'front received event at {datetime.datetime.now()}',
        )
    except Exception as e:
        print(e)
        return {'statusCode': 500}
    print('sent')


def handle_worker(event, context):
    print(f'worker received event: {event["body"]}')
