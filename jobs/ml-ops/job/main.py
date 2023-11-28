import sys, os, pickle, boto3
import ml_training as ml
from dotenv import load_dotenv

def main() -> int:
    """
    This function trains a classifier on data pulled from a data store.
    It uploads training/test artifacts into object storage.
    """

    load_dotenv(dotenv_path='./.env')

    s3 = boto3.resource(
        's3',
        region_name=os.getenv('SCW_REGION'),
        use_ssl=True,
        endpoint_url='https://s3.{region}.scw.cloud'.format(region=os.getenv('SCW_REGION')),
        aws_access_key_id=os.getenv('SCW_ACCESS_KEY'),
        aws_secret_access_key=os.getenv('SCW_SECRET_KEY'),
    )

    # download data locally from data store
    data_store = s3.Bucket(name=os.getenv('SCW_DATA_STORE')) # type: ignore
    data_store.download_file(os.getenv('DATA_FILE_NAME'), './data/'+os.getenv('DATA_FILE_NAME', ''))
    data = ml.load_data('./data/'+os.getenv('DATA_FILE_NAME', ''))
    cleaned_data = ml.clean_data(data)
    transformed_data = ml.transform_data(cleaned_data)

    X_train, X_test, y_train, y_test = ml.split_to_train_test_data(transformed_data)
    X_train, y_train = ml.over_sample_target_class(X_train, y_train)

    # train and upload classifier to model registry
    classifier, _ = ml.tune_classifier(X_train, y_train)
    pickle.dump(classifier, open(os.getenv('MODEL_FILE_NAME',''), 'wb'))
    model_registry = s3.Bucket(name=os.getenv('SCW_MODEL_REGISTRY')) # type: ignore
    model_registry.upload_file(Filename='/ml-job/'+os.getenv('MODEL_FILE_NAME', ''), Key=os.getenv('MODEL_FILE_NAME'))

    # compute performance on test data
    y_pred = ml.predict_on_test_data(classifier, X_test)
    y_pred_prob = ml.predict_prob_on_test_data(classifier,X_test)
    test_metrics = ml.compute_performance_metrics(y_test, y_pred, y_pred_prob)
    pickle.dump(test_metrics, open('performance_metrics.pkl', 'wb'))
    performance_monitor = s3.Bucket(name=os.getenv('SCW_PERF_MONITOR')) # type: ignore
    performance_monitor.upload_file(Filename='/ml-job/performance_metrics.pkl', Key='performance_metrics.pkl')

    # save roc_auc plot
    ml.save_roc_plot(classifier, X_test, y_test)
    performance_monitor.upload_file(Filename='/ml-job/roc_auc_curve.png', Key='roc_auc_curve.png')

    # save confusion matrix
    ml.save_confusion_matrix_plot(classifier, X_test, y_test)
    performance_monitor.upload_file(Filename='/ml-job/confusion_matrix.png', Key='confusion_matrix.png')

    return 0

if __name__ == '__main__':
    sys.exit(main())