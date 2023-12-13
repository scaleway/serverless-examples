from fastapi import FastAPI
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import RocCurveDisplay
import pickle
import boto3
import pandas
import os

import data

classifier = RandomForestClassifier()

app = FastAPI()

MODEL_FILE = "classifier.pkl"


class ClassifierLoader:
    _classifier = None

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

            s3.download_file(bucket_name, MODEL_FILE, MODEL_FILE)

            with open(MODEL_FILE, "rb") as fh:
                cls._classifier = pickle.load(fh)

        return cls._classifier


@app.post("/load")
def load():
    """Reloads classifier from model registry bucket"""

    ClassifierLoader.load(force=True)
    return {"message": "model loaded successfully"}


@app.post("/inference")
def classify(profile: data.ClientProfile):
    """Predicts class given client profile"""

    cleaned_data = data.clean_profile(profile)
    data_point_processed = data.transform_data(cleaned_data)

    # Lazy-loads classifer from S3
    classifier = ClassifierLoader.load()
    prediction = classifier.predict(data_point_processed)

    return {"predicted_class": int(prediction)}
