import random


# This function receives an event from SQS and then randomly succeeds or fails. According to
# the result, the event will be considered delivered or a retry will happen.
def handle(event, context):
    print(f"Received body: {event['body']}")

    # In the case of trigger activated handlers, even if possible, it doesn't make much sense to return anything
    # on success. However, returning an error code will force a redelivery of the event after a delay.
    if random.choice([True, False]):
        print("Success")
        return
    else:
        print("Error")
        return {"statusCode": 500}
