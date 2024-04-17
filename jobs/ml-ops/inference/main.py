from fastapi import FastAPI
from sklearn.ensemble import RandomForestClassifier
from loader import ClassifierLoader
import data

classifier = RandomForestClassifier()

app = FastAPI()

# this endpoint is used by cron trigger to load model from S3
@app.get("/")
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
