import os

from fastapi import FastAPI
from llama_cpp import Llama
from pydantic import BaseModel


class Message(BaseModel):
    content: str


MODEL_FILE_NAME = os.environ["MODEL_FILE_NAME"]

app = FastAPI()

print("loading model starts", flush=True)

llm = Llama(model_path=MODEL_FILE_NAME)

print("loading model successfully ends", flush=True)


@app.get("/")
def hello():
    """Get info of inference server"""

    return {
        "message": "Hello, this is the inference server! Serving model {model_name}".format(
            model_name=MODEL_FILE_NAME
        )
    }


@app.post("/")
def infer(message: Message):
    """Post a message and receive a response from inference server"""

    print("inference endpoint is called", flush=True)

    output = llm(prompt=message.content, max_tokens=200)

    print("output is successfully inferred", flush=True)

    print(output, flush=True)

    return output
