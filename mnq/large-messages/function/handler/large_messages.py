import os
import boto3
from botocore.exceptions import ClientError
import img2pdf
from PIL import Image

endpoint_url = os.getenv("ENDPOINT_URL")
bucket_region = os.getenv("BUCKET_REGION")
bucket_name = os.getenv("BUCKET_NAME")
access_key_id = os.getenv("ACCESS_KEY_ID")
secret_access_key = os.getenv("SECRET_ACCESS_KEY")


def convert_img_to_pdf(img_path, pdf_path):
    image = Image.open(img_path)
    pdf_bytes = img2pdf.convert(image.filename)
    file = open(pdf_path, "wb")
    file.write(pdf_bytes)
    image.close()
    file.close()
    print("Successfully made pdf file")


def handle(event, context):
    input_file = event["body"]
    output_file = os.path.splitext(input_file)[0] + ".pdf"
    s3 = boto3.client(
        "s3",
        endpoint_url=endpoint_url,
        region_name=bucket_region,
        aws_access_key_id=access_key_id,
        aws_secret_access_key=secret_access_key,
    )

    try:
        s3.download_file(bucket_name, input_file, input_file)
        print("Object " + input_file + " downloaded")
        convert_img_to_pdf(input_file, output_file)
        s3.upload_file(output_file, bucket_name, output_file)
        print("Object " + input_file + " uploaded")
    except ClientError as e:
        print(e)
        return {
            "body": {"message": e.response["Error"]["Message"]},
            "statusCode": e.response["Error"]["Code"],
        }

    return {"body": {"message": "Converted in PDF sucessfully"}, "statusCode": 200}
