from fastapi import FastAPI
from sklearn.ensemble import RandomForestClassifier
import data_processing as process
import pickle, boto3, pandas, os

classifier = RandomForestClassifier()

app = FastAPI()


@app.get("/load_classifier")
def load_classifier():
    """(Re)loads classifier from model registry bucket"""

    s3 = boto3.resource(
        "s3",
        region_name=os.environ["MAIN_REGION"],
        use_ssl=True,
        endpoint_url=f'https://s3.{os.environ["MAIN_REGION"]}.scw.cloud',
        aws_access_key_id=os.environ["SCW_ACCESS_KEY"],
        aws_secret_access_key=os.environ["SCW_SECRET_KEY"],
    )

    bucket = s3.Bucket(name=os.environ["MODEL_REGISTRY"])  # type: ignore
    bucket.download_file(os.environ["MODEL_FILE"], os.environ["MODEL_FILE"])

    global classifier
    classifier = pickle.load(open(os.environ["MODEL_FILE"], "rb"))

    return {"message": "model loaded successfully"}


@app.post("/inference")
def classify(data: process.ClientProfile):
    """Predicts class given client profile"""

    data_point_json = data.model_dump()
    data_point_pd = pandas.DataFrame(index=[0], data=data_point_json)
    data_point_processed = process.transform_data(process.clean_data(data_point_pd))
    global classifier
    prediction = classifier.predict(data_point_processed)

    return {"predicted_class": int(prediction)}
