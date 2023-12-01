import pandas as pd
import os
import pickle
import boto3
import training as ml
from sklearn.metrics import RocCurveDisplay
from sklearn.metrics import ConfusionMatrixDisplay

DATA_FILE_NAME = "bank-additional-full.csv"

MODEL_FILE = "classifier.pkl"
PERF_FILE = "performance.pkl"
ROC_AUC_FILE = "roc_auc.png"
CONFUSION_MATRIX_FILE = "confusion_matrix.png"


def main() -> int:
    """
    Trains a classifier on data pulled from a data store.
    Uploads training/test artifacts into artifact data stores.
    """

    access_key = os.environ["SCW_ACCESS_KEY"]
    secret_key = os.environ["SCW_SECRET_KEY"]
    region_name = os.environ["SCW_REGION"]

    bucket_name = os.environ["S3_BUCKET_NAME"]
    s3_url = os.environ["S3_URL"]

    s3 = boto3.client(
        "s3",
        region_name=region_name,
        endpoint_url=s3_url,
        aws_access_key_id=access_key,
        aws_secret_access_key=secret_key,
    )

    # Download data
    print(f"Downloading data from {s3_url}/{bucket_name}/{DATA_FILE_NAME}")
    s3.download_file(bucket_name, DATA_FILE_NAME, DATA_FILE_NAME)
    data = pd.read_csv(DATA_FILE_NAME, sep=";")

    # Clean and transform data
    cleaned_data = data.dropna()
    transformed_data = ml.transform_data(cleaned_data)

    # Split train and test
    x_train, x_test, y_train, y_test = ml.split_to_train_test_data(transformed_data)
    x_train, y_train = ml.over_sample_target_class(x_train, y_train)

    # Train and upload classifier to s3
    classifier, _ = ml.tune_classifier(x_train, y_train)

    with open(MODEL_FILE, "wb") as fh:
        pickle.dump(classifier, fh)

    print(f"Uploading model to {s3_url}/{bucket_name}/{MODEL_FILE}")
    s3.upload_file(MODEL_FILE, bucket_name, MODEL_FILE)

    # Compute performance on test data
    y_pred = classifier.predict(x_test)
    y_pred_prob = classifier.predict_proba(x_test)
    test_metrics = ml.compute_performance_metrics(y_test, y_pred, y_pred_prob)

    with open(PERF_FILE, "wb") as fh:
        pickle.dump(test_metrics, fh)

    print(f"Uploading perf to {s3_url}/{bucket_name}/{PERF_FILE}")
    s3.upload_file(PERF_FILE, bucket_name, PERF_FILE)

    # save roc_auc plot
    print(f"Uploading ROC curve to {s3_url}/{bucket_name}/{ROC_AUC_FILE}")
    display = RocCurveDisplay.from_estimator(classifier, x_test, y_test)
    display.figure_.savefig(ROC_AUC_FILE)
    s3.upload_file(ROC_AUC_FILE, bucket_name, ROC_AUC_FILE)

    # save confusion matrix
    print(
        f"Uploading confusion matrix to {s3_url}/{bucket_name}/{CONFUSION_MATRIX_FILE}"
    )
    display = ConfusionMatrixDisplay.from_estimator(classifier, x_test, y_test)
    display.figure_.savefig(CONFUSION_MATRIX_FILE)
    s3.upload_file(CONFUSION_MATRIX_FILE, bucket_name, CONFUSION_MATRIX_FILE)


if __name__ == "__main__":
    main()
