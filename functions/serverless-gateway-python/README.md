# Using Serverless Gateway and Serverless APIs

This example shows how to serve a Serverless API composed of many functions through our serverless gateway.

## Requirements

This example uses:
* [Python API Framework](https://github.com/scaleway/serverless-api-project) to deploy the functions.
* [Serverless Gateway](https://github.com/scaleway/serverless-gateway) to deploy the a serverless gateway as a container.

## Gateway set-up

Deploy the serverless gateway as a container following the Serverless Gateway project instructions.

Make sure to export your gateway base URL (`GATEWAY_URL`) and an authentication token (`$GATEWAY_TOKEN`)

## Running 

### Deploy your serverless API

Deploy your functions using:

```
pip install -r requirements.txt
scw_serverless deploy app.py
```

### Add your functions to gateway endpoints

You can add your function to the gateway with a chosen relative path using:
```
sh scripts/add_function_to_gateway.sh http://<function_domain_name> /<chosen_relative_path>
```

### Verify that your function is added to the gateway:

You can use:
```
sh list_gateway_endpoints.sh
```

### Call your functions through gateway

You can use:
```
curl scripts/http://${GATEWAY_URL}/<chosen_relative_path>
```

### Delete your function through gateway

You can use:
```
sh scripts/delete_function_to_gateway.sh http://<function_domain_name> /<chosen_relative_path>
```

