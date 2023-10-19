from flask import Flask, make_response

app = Flask(__name__)

HTTP_METHODS = ["GET", "POST"]


@app.route("/", methods=HTTP_METHODS)
def root():
    response = make_response("Hello from container")
    return response


if __name__ == "__main__":
    app.run(debug=True, host="0.0.0.0", port="8080")
