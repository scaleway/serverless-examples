import os
from urllib import request,parse,error
import json

auth_token=os.environ['X-AUTH-TOKEN']

def handle(event, context):
    # Get information from cron
    event_body=eval(event["body"])
    zone=event_body["zone"]
    server_id=event_body["server_id"]
    action=event_body["action"] # action should be "poweron" or "poweroff"

    # Create request
    url=f"https://api.scaleway.com/instance/v1/zones/{zone}/servers/{server_id}/action"
    data=json.dumps({"action":action}).encode('ascii')
    req = request.Request(url, data=data,  method="POST")
    req.add_header('Content-Type', 'application/json')
    req.add_header('X-Auth-Token',auth_token)

    # Sending request to Instance API
    try:
        res=request.urlopen(req).read().decode()
    except error.HTTPError as e:
        res=e.read().decode()

    return {
        "body": json.loads(res),
        "statusCode": 200,
    }
