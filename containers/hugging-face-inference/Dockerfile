FROM python:3.12-slim-bookworm

ARG MODEL_DOWNLOAD_SOURCE

RUN apt-get update && apt-get install -y wget

WORKDIR /app

RUN pip install --upgrade pip
COPY requirements.txt .
RUN pip install -r requirements.txt

RUN pip install llama-cpp-python==0.2.62 \
    --no-cache-dir \
    --extra-index-url https://abetlen.github.io/llama-cpp-python/whl/cpu

RUN wget $MODEL_DOWNLOAD_SOURCE

COPY . .

CMD ["uvicorn", "main:app", "--proxy-headers", "--host", "0.0.0.0", "--port", "80"]
