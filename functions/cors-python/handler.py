def handle(event, context):

    print(event)
    print(context)

    return {
        "body": "This is checking CORS",
        "headers": {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Headers": "*",
            "Access-Control-Allow-Methods": "*",
            "Content-Type": ["text/plain"],
        },
    }

if __name__ == "__main__":
    from scaleway_functions_python import local
    local.serve_handler(handle)