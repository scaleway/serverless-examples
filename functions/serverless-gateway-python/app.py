from  scw_serverless import Serverless

app = Serverless("serverless-api")

@app.get(url="/func-a")
def func_a(_event, _context):
    return "Hello from function A"

@app.get(url="/func-b")
def func_b(_event, _context):
    return "Hello from function B"

@app.get(url="/func-c")
def func_c(_event, _context):
    return "Hello from function C"

if __name__ == "__main__":
    from scaleway_functions_python import local

    server = local.LocalFunctionServer()
    server.add_handler(func_a)
    server.add_handler(func_b)
    server.add_handler(func_c)
    server.serve(port=8080)
