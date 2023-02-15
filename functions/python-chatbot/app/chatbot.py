from chatterbot import ChatBot
import json
import os

script_dir = os.path.dirname(__file__)

# Instantiate chatbot
abs_db_path = os.path.join(script_dir, os.environ['database_model'])
chatbot = ChatBot(
    'English Bot',
    database_uri="sqlite:///" + abs_db_path
)


def generate_response(body, status_code, content_type):
    return {
        "body": body,
        "statusCode": status_code,
        "headers": {
            "Content-Type": [content_type]
        }
    }


def render_html():
    abs_file_path = os.path.join(script_dir, "templates/index.html")
    return generate_response(open(abs_file_path, "r").read(), 200, "text/html")


def chatbot_response(event):
    try:
        body = json.loads(event['body'])
        message = body["message"]
    except (json.decoder.JSONDecodeError, KeyError):
        return generate_response('The parameter "message" is required.', 400, "text/plain")

    # Return chatterbot response
    response = chatbot.get_response(message).text
    return generate_response({"response": response}, 200, "application/json")


def handle(event, context):
    method = event['httpMethod']
    match method:
        case'GET':
            return render_html()
        case 'POST':
            return chatbot_response(event)
        case _:
            return generate_response('GET and POST methods only are allowed', 405, "text/plain")
