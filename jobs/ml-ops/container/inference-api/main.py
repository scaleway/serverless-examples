from fastapi import FastAPI
from dotenv import load_dotenv
import data_processing as process
import pickle, boto3, pandas, os

app = FastAPI()

load_dotenv(dotenv_path='./.env')

s3 = boto3.resource(
    's3',
    region_name=os.getenv('SCW_REGION'),
    use_ssl=True,
    endpoint_url='https://s3.{region}.scw.cloud'.format(region=os.getenv('SCW_REGION')),
    aws_access_key_id=os.getenv('SCW_ACCESS_KEY'),
    aws_secret_access_key=os.getenv('SCW_SECRET_KEY'),
)

bucket = s3.Bucket(name=os.getenv('SCW_MODEL_REGISTRY')) # type: ignore
bucket.download_file(os.getenv('MODEL_FILE_NAME'), os.getenv('MODEL_FILE_NAME'))

classifier = pickle.load(open(os.getenv('MODEL_FILE_NAME', ''),'rb'))

@app.post('/inference')
def classify(data:process.ClientProfile):
    """
    This function predicts class given client profile
    """
    
    data_point_json = data.model_dump()
    data_point_pd = pandas.DataFrame(index=[0], data=data_point_json)
    data_point_processed = process.transform_data(process.clean_data(data_point_pd))
    prediction = classifier.predict(data_point_processed)

    return {'predicted_class': int(prediction)}