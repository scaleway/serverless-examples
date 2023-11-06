from flask import Flask
import os

app = Flask(__name__)


@app.route("/")
def hello():
    return "Hello world from the Python container!"


if __name__ == "__main__":
    port_env = os.getenv("PORT", 8080)
    port = int(port_env)
    app.run(debug=True, host="0.0.0.0", port=port)
