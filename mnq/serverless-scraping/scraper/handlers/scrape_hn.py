import requests
import boto3
import json 
import os
from datetime import datetime, timedelta
from bs4 import BeautifulSoup

HN_URL = "https://news.ycombinator.com/newest"
SCW_SQS_URL = "https://sqs.mnq.fr-par.scaleway.com"

queue_url = os.getenv('QUEUE_URL') 
sqs_access_key = os.getenv('SQS_ACCESS_KEY')
sqs_secret_access_key = os.getenv('SQS_SECRET_ACCESS_KEY')

def scrape_and_push():
    """
    Scrape the HN website for articles published in the last 15 minutes, and push infos on the SQS queue
    """
    page = requests.get(HN_URL)
    html_doc = page.content

    soup = BeautifulSoup(html_doc, 'html.parser')

    # On hn news page there are exactly 30 articles, for each of them a `titleline` and a `age` span are present
    titlelines = soup.find_all(class_="titleline")
    ages = soup.find_all(class_="age")

    sqs = boto3.client('sqs', endpoint_url=SCW_SQS_URL, aws_access_key_id=sqs_access_key, aws_secret_access_key=sqs_secret_access_key, region_name='fr-par')

    for age, titleline in zip(ages, titlelines):
        time_str = age["title"]
        time = datetime.strptime(time_str, "%Y-%m-%dT%H:%M:%S")
        # Check if article was published in the last 15 minutes
        if datetime.utcnow() - time > timedelta(minutes=15):
            continue

        body = json.dumps({'url': titleline.a["href"], 'title': titleline.a.get_text()})
        response = sqs.send_message(QueueUrl=queue_url, MessageBody=body)
        
    return page.status_code

def handle(event, context):
    try:
        status = scrape_and_push()
        return {'statusCode': status, 'headers': {'content': 'text/plain'}}
    except Exception as e:
        print(e)
        return {'statusCode': 500, 'body': str(e)}

if __name__ == "__main__":
    handle(None, None)