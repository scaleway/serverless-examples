def handle(event, context):
    # The content of the NATS message is passed in the body
    msg_body = event.get("body")
    print(f"Received message: {msg_body}")

    return "Hello NATS!"
