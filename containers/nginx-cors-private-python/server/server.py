from flask import Flask, send_file, make_response, request
from datetime import datetime
from os.path import exists

app = Flask(__name__)

AUTH_HEADER = "X-Auth-Token"

DUMMY_IMAGE = "/images/dummy.png"
MIME_TYPE = "image/png"

HTTP_METHODS = ["GET", "POST"]


def add_cors_headers(response):
    response.headers.add("Access-Control-Allow-Origin", "*")
    response.headers.add("Access-Control-Allow-Headers", "*")
    response.headers.add("Access-Control-Allow-Methods", "*")

    return response


def check_auth_header():
    header = request.headers.get(AUTH_HEADER)
    if not header:
        app.logger.warn("No auth header %s", AUTH_HEADER)


@app.route("/", methods=HTTP_METHODS)
def root():
    now = datetime.now()
    req_time = now.strftime("%H:%M:%S")
    check_auth_header()

    response = make_response(f"Hello at {req_time}")
    return add_cors_headers(response)


@app.route("/img/<filename>", methods=HTTP_METHODS)
def img(filename):
    img_path = f"/images/{filename}"
    check_auth_header()

    # Check if image exists
    img_path = img_path if exists(img_path) else DUMMY_IMAGE

    app.logger.info("Serving image at %s", img_path)

    response = send_file(img_path, mimetype=MIME_TYPE)
    return add_cors_headers(response)


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port="8080")
