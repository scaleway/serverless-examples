import psycopg2
from psycopg2 import Error
import os
import logging

PG_HOST=os.environ.get('PG_HOST')
PG_USER=os.environ.get('PG_USER')
PG_DATABASE=os.environ.get('PG_DATABASE')
PG_PASSWORD=os.environ.get('PG_PASSWORD')
PG_PORT=os.environ.get('PG_PORT')
PG_SSL_ROOT_CERT=os.environ.get('PG_SSL_ROOT_CERT')

def handle(event, context):

    try:
        connection = psycopg2.connect(
            database=PG_DATABASE,
            user=PG_USER,
            host=PG_HOST,
            password=PG_PASSWORD,
            port=PG_PORT,
            sslmode="require",
            sslrootcert=PG_SSL_ROOT_CERT,
            )

    except (Exception, Error) as error:
        logging.error("Error while connecting to PostgreSQL database", error)
        return {
                "statusCode": 500,
                "body": {
                    "message": "Error while connecting to PostgreSQL database, check function logs for more information"
                }
        }

    try:
        cursor = connection.cursor()
        logging.info("Connected to Database")
        cursor.execute("SELECT * FROM table LIMIT 10")
        record = cursor.fetchone()
        logging.info("Successfully fetched data")
        cursor.close()
        return {
                "statusCode": 200,
                "body": {
                    "message": record
                }
        }

    except (Exception, Error) as error:
        logging.error("Error while interacting with PostgreSQL", error)
        return {
                "statusCode": 500,
                "body": {
                    "message": "Error while getting information from PostgreSQL database, check function logs for more information"
                }
        }

if __name__ == "__main__":
    from scaleway_functions_python import local
    local.serve_handler(handle)
