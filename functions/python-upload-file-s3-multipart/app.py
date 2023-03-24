from typing import TYPE_CHECKING
import logging
import os

from scw_serverless import Serverless
if TYPE_CHECKING:
    from scaleway_functions_python.framework.v1.hints import Context, Event, Response

import boto3
from streaming_form_data import StreamingFormDataParser
from streaming_form_data.targets import ValueTarget

SCW_ACCESS_KEY = os.environ["SCW_ACCESS_KEY"]
SCW_SECRET_KEY = os.environ["SCW_SECRET_KEY"]
BUCKET_NAME = os.environ["BUCKET_NAME"]

# Files will be uploaded to cold storage
# See: https://www.scaleway.com/en/glacier-cold-storage/
STORAGE_CLASS = "GLACIER"

app = Serverless(
    "s3-utilities",
    secret={
        "SCW_ACCESS_KEY": SCW_ACCESS_KEY,
        "SCW_SECRET_KEY": SCW_SECRET_KEY,
    },
    env={
        "BUCKET_NAME": BUCKET_NAME,
        "PYTHONUNBUFFERED": "1",
    },
)

s3 = boto3.resource(
    "s3",
    region_name="fr-par",
    use_ssl=True,
    endpoint_url="https://s3.fr-par.scw.cloud",
    aws_access_key_id=SCW_ACCESS_KEY,
    aws_secret_access_key=SCW_SECRET_KEY,
)

bucket = s3.Bucket(BUCKET_NAME)

logging.basicConfig(level=logging.INFO)


@app.func()
def upload(event: "Event", _context: "Context") -> "Response":
    """Upload form data to S3 Glacier."""

    headers = event["headers"]
    parser = StreamingFormDataParser(headers=headers)

    # Defined a target to handle files uploaded with the "file" key
    target = ValueTarget()
    parser.register("file", target)

    body: str = event["body"]
    parser.data_received(body.encode("utf-8"))

    if not (len(target.value) > 0 and target.multipart_filename):
        return {"statusCode": 400}

    name = target.multipart_filename

    logging.info("Uploading file %s to Glacier on %s", name, bucket.name)
    bucket.put_object(Key=name, Body=target.value, StorageClass=STORAGE_CLASS)

    return {"statusCode": 200, "body": f"Successfully uploaded {name} to bucket!"}


if __name__ == "__main__":
    from scaleway_functions_python import local

    local.serve_handler(upload)
