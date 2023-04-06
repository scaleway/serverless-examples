from  scw_serverless import Serverless

app = Serverless("serverless-api")

@app.func(relative_url="/func-a")
def func_a(_event, _context):
    return "Hello from function A"

@app.func(relative_url="/func-b")
def func_b(_event, _context):
    return "Hello from function B"

@app.func(relative_url="/func-c")
def func_c(_event, _context):
    return "Hello from function C"