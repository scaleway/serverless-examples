from flask import Flask, request
import logging
import os

import boto3

REGION = "fr-par"
S3_URL = "https://s3.fr-par.scw.cloud"

SCW_ACCESS_KEY = os.environ["ACCESS_KEY"]
SCW_SECRET_KEY = os.environ["SECRET_KEY"]
BUCKET_NAME = os.environ["BUCKET_NAME"]

logging.basicConfig(level=logging.INFO)

app = Flask(__name__)


@app.route("/", methods=["POST", "GET"])
def main():
    s3 = boto3.client(
        "s3",
        region_name=REGION,
        use_ssl=True,
        endpoint_url=S3_URL,
        aws_access_key_id=SCW_ACCESS_KEY,
        aws_secret_access_key=SCW_SECRET_KEY,
    )

    uploaded_file = request.files["file"]
    file_body = uploaded_file.read()

    logging.info(f"Uploading to {BUCKET_NAME}/{uploaded_file.filename}")

    s3.put_object(
        Key=uploaded_file.filename, Bucket=BUCKET_NAME, Body=file_body
    )

    return {
        "statusCode": 200,
        "body": f"Successfully uploaded {uploaded_file.filename} to bucket!",
    }


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port=8080)
