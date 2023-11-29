from dotenv import load_dotenv
import os, sys, boto3


def main() -> int:
    """Uploads a local CSV file to a target S3 bucket"""

    load_dotenv(dotenv_path="./.env")

    s3 = boto3.resource(
        "s3",
        region_name=os.environ["SCW_REGION"],
        use_ssl=True,
        endpoint_url=f'https://s3.{os.environ["SCW_REGION"]}.scw.cloud',
        aws_access_key_id=os.environ["SCW_ACCESS_KEY"],
        aws_secret_access_key=os.environ["SCW_SECRET_KEY"],
    )

    bucket = s3.Bucket(name=os.environ["SCW_DATA_STORE"])  # type: ignore
    bucket.upload_file(
        Filename="./data/" + os.environ["DATA_FILE_NAME"],
        Key=os.environ["DATA_FILE_NAME"],
    )

    return 0


if __name__ == "__main__":
    sys.exit(main())
