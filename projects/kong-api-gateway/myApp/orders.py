from jinja2 import Environment, DictLoader, select_autoescape

def order(event, context):
    index_html = """
    <html>
      <body>
        orders
      </body>
    </html>
    """
    env = Environment(
         loader=DictLoader({'index.html': index_html}),
         autoescape=select_autoescape()
    )
    template = env.get_template("index.html")
    return {
        "body": template.render(),
        "statusCode": 200,
    }

