FROM python:3-alpine
WORKDIR /usr/src/app

COPY requirements.txt .
RUN pip install -qr requirements.txt
COPY server.py .

WORKDIR /usr/src/app
CMD ["python3", "./server.py"]
