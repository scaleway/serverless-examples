import psycopg2
from psycopg2 import Error
import os
import logging

PG_HOST=os.getenv('PG_HOST')
PG_USER=os.getenv('PG_USER')
PG_DATABASE=os.getenv('PG_DATABASE')
PG_PASSWORD=os.getenv('PG_PASSWORD')
PG_PORT=os.getenv('PG_PORT')

def handle(event, context):

    try:
        connection = psycopg2.connect(
            database=PG_DATABASE,
            user=PG_USER,
            host=PG_HOST,
            password=PG_PASSWORD,
            port=PG_PORT,
            sslmode="disable",
        )
        logging.info("Connected to Database")


    except (Exception, Error) as error:
        logging.error("Error while connecting to PostgreSQL database", error)
        return {
                "statusCode": 500,
                "body": {
                    "message": "Error while connecting to PostgreSQL database, check function logs for more information"
                }
        }

    try:
        sql_create_table = """CREATE TABLE dummy_table (
            field_1 int,
            field_2 varchar(20)
        );"""

        sql_insert_new_lines = """INSERT INTO dummy_table (field_1, field_2) 
            VALUES
            (10, 'dummy_data_1'),
            (20, 'dummy_data_2');
        """

        with connection.cursor() as cursor:
            logging.info("Successfully created database cursor")
            cursor.execute(sql_create_table)
            cursor.execute(sql_insert_new_lines)
            connection.commit()
            cursor.close()
            logging.info("Successfully executed SQL queries")
        
        connection.close()
        logging.info("Database connection is closed")
        
        return {
                "statusCode": 200,
                "body": {
                    "message": "successful sql queries"
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
