import requests


def handle(event, context):
    resp = requests.get("https://www.scaleway.com")
    return {
        "body": f"Response status: {resp.status_code}",
        "headers": {
            "Content-Type": ["text/plain"],
        },
    }
