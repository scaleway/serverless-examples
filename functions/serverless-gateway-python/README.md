# Using Serverless Gateway and Serverless APIs

This example shows how to serve a Serverless API composed of many functions through our serverless gateway.

## Requirements

This example uses:
* [Python API Framework](https://github.com/scaleway/serverless-api-project) to deploy the functions.
* [Serverless Gateway](https://github.com/scaleway/serverless-gateway) to deploy a serverless gateway container.

## Gateway set-up

Deploy the serverless gateway as a container following the Serverless Gateway project instructions.

*Make sure to export your gateway base URL (`GATEWAY_HOST`) and an authentication token (`TOKEN`)*

## Running 

### Deploy your serverless API

Deploy your functions and add them automatically as endpoints to your serverless gateway using:

```
pip install -r requirements.txt
scw-serverless deploy app.py --gateway-url https://${GATEWAY_HOST} --gateway-api-key ${TOKEN}
```

### Check your endpoint has been added

You can use:
```
./scripts/list_gateway_endpoints.sh
```

### Call your function via its route

You can use:
```
curl https://${GATEWAY_HOST}/<chosen_relative_path>
```

### Delete your endpoint

You can use:
```
./scripts/delete_function_to_gateway.sh https://<function_domain_name> /<chosen_relative_path>
```

