from datetime import datetime
import json
import os
from pathlib import Path

import redis
import requests

WEATHER_API_URL = "https://api.open-meteo.com/v1/forecast"
PARIS_LATITUDE = 48.856614
PARIS_LONGITUDE = 2.3522219

REDIS_URL = os.environ["REDIS_URL"]
REDIS_USER = os.environ["REDIS_USER"]
REDIS_PASSWORD = os.environ["REDIS_PASSWORD"]

REDIS_CERT_PATH = "/etc/ssl/certs/redis/redis_certificate.crt"


def handle(event, context):
    cert_path = Path(REDIS_CERT_PATH)
    if not cert_path.exists():
        cert_path.parent.mkdir(exist_ok=True, parents=True)
        cert_path.write_text(os.environ["REDIS_CERT"])
        cert_path.chmod(0o600)

    red = redis.Redis(
        host=REDIS_URL,
        username=REDIS_USER,
        password=REDIS_PASSWORD,
        ssl=True,
        ssl_ca_certs=str(cert_path.resolve()),
    )
    red.ping()

    resp = requests.get(
        url=WEATHER_API_URL,
        params={
            "longitude": PARIS_LONGITUDE,
            "latitude": PARIS_LATITUDE,
            "current_weather": True,
            "hourly": ["temperature_2m"],
        },
    )
    resp.raise_for_status()
    data = resp.json()

    hourly_temperature = {
        k.split("T")[1]: v
        for (k, v) in zip(data["hourly"]["time"], data["hourly"]["temperature_2m"])
    }
    today = datetime.today().strftime('%Y-%m-%d')
    print(f"Hourly temperature for {today}: {hourly_temperature}")

    red.set(today, json.dumps(hourly_temperature))

    return {
        "body": "Successfully retrieved temperatures for today in Paris",
        "headers": {
            "Content-Type": ["text/plain"],
        }
    }
