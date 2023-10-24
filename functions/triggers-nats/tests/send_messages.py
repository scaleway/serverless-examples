import asyncio
import os

import nats

NATS_ENDPOINT = os.environ.get("NATS_ENDPOINT")
NATS_CREDS_PATH = "../nats-creds"

NATS_TOPIC = "triggers-nats-topic"


async def main():
    print(f"Connecting to NATS account at {NATS_ENDPOINT}")

    # Connect to NATS
    nc = await nats.connect(servers=[NATS_ENDPOINT], user_credentials=NATS_CREDS_PATH)

    # Publish a message
    await nc.publish(NATS_TOPIC, b"Hello NATS triggers!")

    # Close NATS connection
    await nc.close()


if __name__ == "__main__":
    asyncio.run(main())
