FROM python:3.12-slim-bookworm

WORKDIR /ml-job/

COPY . .

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

CMD [ "python", "main.py" ]
