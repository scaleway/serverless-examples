import pickle
import boto3
import os

class ClassifierLoader:
    _classifier = None
    _classifier_version = ""

    @classmethod
    def load(cls, force=False):
        if force or cls._classifier is None:
            access_key = os.environ["ACCESS_KEY"]
            secret_key = os.environ["SECRET_KEY"]
            region_name = os.environ["REGION"]

            bucket_name = os.environ["S3_BUCKET_NAME"]
            s3_url = os.environ["S3_URL"]

            s3 = boto3.client(
                "s3",
                region_name=region_name,
                endpoint_url=s3_url,
                aws_access_key_id=access_key,
                aws_secret_access_key=secret_key,
            )

            # get model file with the latest version
            bucket_objects = s3.list_objects(Bucket=bucket_name)
            get_last_modified = lambda object: int(object['LastModified'].strftime('%s'))
            model_objects = [model_object for model_object in bucket_objects['Contents'] if "classifier" in model_object['Key']]
            latest_model_file = [object['Key'] for object in sorted(model_objects, key=get_last_modified)][0]

            s3.download_file(bucket_name, latest_model_file, latest_model_file)

            with open(latest_model_file, "rb") as fh:
                cls._classifier = pickle.load(fh)
                cls._classifier_version = latest_model_file[11:-4]

            print('Successfully loaded model file: {latest_model_file}'.format(latest_model_file=latest_model_file), flush=True)

        return cls._classifier

    @classmethod
    def model_version(cls):
        return cls._classifier_version