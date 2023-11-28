FROM python:3.10.12

WORKDIR /ml-job/

COPY . .

RUN pip install --upgrade pip
RUN pip install -r requirements.txt

CMD [ "python", "main.py" ]
