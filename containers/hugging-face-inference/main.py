from fastapi import FastAPI
from llama_cpp import Llama
import os
import prompt

MODEL_FILE_NAME=os.environ["MODEL_FILE_NAME"]

app = FastAPI()

print("loading model from memory starts", flush=True)

llm = Llama(model_path=MODEL_FILE_NAME)

print("loading model from memory successfully ends", flush=True)

@app.get("/")
def hello():
    """Get Inference Server Info"""

    return {
         "message": "Hello, this is the inference server! Serving model {model_name}"
         .format(model_name=MODEL_FILE_NAME)
    }

@app.post("/")
def infer(prompt: prompt.Prompt):

    print("inference endpoint is called", flush=True)

    output = llm(prompt=prompt.message, max_tokens=200)

    print("output is successfully inferred", flush=True)

    print(output, flush=True)

    return output
