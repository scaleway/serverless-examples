FROM python:3.12-slim-bookworm

WORKDIR /app

RUN pip install --upgrade pip
COPY requirements.txt .
RUN pip install -r requirements.txt

COPY . .
CMD [ "python", "main.py" ]
