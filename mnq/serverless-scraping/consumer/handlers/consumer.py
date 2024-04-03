import json 
import os
import pg8000.native
import requests
from bs4 import BeautifulSoup

db_host = os.getenv('DB_HOST')
db_port = os.getenv('DB_PORT')
db_name = os.getenv('DB_NAME')
db_user = os.getenv('DB_USER')
db_password = os.getenv('DB_PASSWORD')

CREATE_TABLE_IF_NOT_EXISTS = """
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    a_count INTEGER NOT NULL,
    h1_count INTEGER NOT NULL, 
    p_count INTEGER NOT NULL
);"""

INSERT_INTO_ARTICLES = """
INSERT INTO articles (title, url, a_count, h1_count, p_count)
VALUES(:title, :url, :a_count, :h1_count, :p_count) RETURNING id
;"""

def scrape_page_for_stats(url):
    """
    Scrape page at given url and return stats about chosen tags
    """
    # articles hosted on hn have a relative url
    if url[:4] == "item":
        url = "https://news.ycombinator.com/" + url

    page = requests.get(url, timeout=15)
    html_doc = page.content
    soup = BeautifulSoup(html_doc, 'html.parser')
    
    tags = ['a', 'h1', 'p']

    return {tag: len(soup.find_all(tag)) for tag in tags}

def scrape_and_save_to_db(event):
    """
    Scrape a page for info and save such infos in db
    """
    body = json.loads(event["body"])

    tags_count = scrape_page_for_stats(body['url'])
    conn = None
    try: 
        conn = pg8000.native.Connection(host=db_host, database=db_name, port=db_port, user=db_user, password=db_password, timeout=15)

        # Where else could we create the table, to avoid manual intervention? 
        conn.run(CREATE_TABLE_IF_NOT_EXISTS)
        conn.run(INSERT_INTO_ARTICLES, title=body['title'], url=body['url'], a_count=tags_count['a'], h1_count=tags_count['h1'], p_count=tags_count['p'])

    finally:
        if conn is not None:
            conn.close()
    return 200

def handle(event, context):
    try:
        status = scrape_and_save_to_db(event)
        return {'statusCode': status, 'headers': {'content': 'text/plain'}}
    except Exception as e:
        print("error", e)
        return {'statusCode': 500, 'body': str(e)}

if __name__ == '__main__':
    handle({'body': json.dumps({'url': 'https://google.com', 'title': 'test url'})}, None)