from  scw_serverless import Serverless

app = Serverless("serverless-api")

@app.func()
def func_a(_event, _context):
    return "Hello from function A"

@app.func()
def func_b(_event, _context):
    return "Hello from function B"

@app.func()
def func_c(_event, _context):
    return "Hello from function C"