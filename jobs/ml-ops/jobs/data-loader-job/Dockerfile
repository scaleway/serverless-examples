FROM python:3.12-slim-bookworm

WORKDIR /data-loader-job

COPY . .

# Get data from public source
RUN apt-get update
RUN apt-get install curl unzip --assume-yes
RUN curl http://archive.ics.uci.edu/static/public/222/bank+marketing.zip -o bank-telemarketing.zip
RUN unzip bank-telemarketing.zip
RUN unzip bank-additional.zip
RUN cp ./bank-additional/bank-additional-full.csv  ./data/
RUN mv data/bank-additional-full.csv data/bank_telemarketing.csv

# Push data to data store
RUN pip install --upgrade pip
RUN pip install -r requirements.txt

CMD ["python", "main.py"]