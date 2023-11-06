FROM python:3.10
WORKDIR /app

RUN pip3 install --upgrade pip
COPY requirements.txt .
RUN pip3 install -r requirements.txt --target .

COPY handler.py .

CMD [ "python3", "./handler.py" ]
