from dotenv import load_dotenv
import os, sys, boto3

def main() -> int:
    """
    This function uploads a local CSV file to a target S3 bucket
    """

    load_dotenv(dotenv_path="./.env")

    s3 = boto3.resource(
    "s3",
    region_name=os.getenv("SCW_REGION"),
    use_ssl=True,
    endpoint_url="https://s3.{region}.scw.cloud".format(region=os.getenv("SCW_REGION")),
    aws_access_key_id=os.getenv("SCW_ACCESS_KEY"),
    aws_secret_access_key=os.getenv("SCW_SECRET_KEY"),
    )

    bucket = s3.Bucket(name=os.getenv("SCW_S3_BUCKET")) # type: ignore
    bucket.upload_file(Filename="./data/"+os.getenv("SOURCE_FILE_NAME",""), Key=os.getenv("SOURCE_FILE_NAME"))

    return 0

if __name__ == '__main__':
    sys.exit(main())