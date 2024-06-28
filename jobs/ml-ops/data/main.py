import os
import urllib.request
import zipfile

import boto3

DATA_DIR = "dataset"

ZIP_URL = "http://archive.ics.uci.edu/static/public/222/bank+marketing.zip"
ZIP_DOWNLOAD_PATH = os.path.join(DATA_DIR, "downloaded.zip")
NESTED_ZIP_PATH = os.path.join(DATA_DIR, "bank-additional.zip")

DATA_FILE = "bank-additional-full.csv"
DATA_CSV_PATH = os.path.join(DATA_DIR, "bank-additional", DATA_FILE)


def main():
    """Pulls file from source, and uploads to a target S3 bucket"""

    # Download the zip
    print(f"Downloading data from {ZIP_URL}")
    os.makedirs(DATA_DIR, exist_ok=True)
    urllib.request.urlretrieve(ZIP_URL, ZIP_DOWNLOAD_PATH)

    # Extract
    with zipfile.ZipFile(ZIP_DOWNLOAD_PATH, "r") as fh:
        fh.extractall(DATA_DIR)

    # Remove original zip
    os.remove(ZIP_DOWNLOAD_PATH)

    # Extract zip within the zip
    with zipfile.ZipFile(NESTED_ZIP_PATH, "r") as fh:
        fh.extractall(DATA_DIR)

    # Remove nested zip
    os.remove(NESTED_ZIP_PATH)

    access_key = os.environ["ACCESS_KEY"]
    secret_key = os.environ["SECRET_KEY"]
    region_name = os.environ["REGION"]

    bucket_name = os.environ["S3_BUCKET_NAME"]
    s3_url = os.environ["S3_URL"]

    print(f"Uploading data to {s3_url}/{bucket_name}")
    s3 = boto3.client(
        "s3",
        region_name=region_name,
        endpoint_url=s3_url,
        aws_access_key_id=access_key,
        aws_secret_access_key=secret_key,
    )

    s3.upload_file(DATA_CSV_PATH, bucket_name, DATA_FILE)


if __name__ == "__main__":
    main()
