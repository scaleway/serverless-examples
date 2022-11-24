scaleway-serverless-apigateway
==============================

This example demonstrate how to setup kong DB-less as API Gateway for Serverless Functions in Scaleway ecosystem.

It uses [Serverless Compose](https://www.serverless.com/framework/docs/guides/compose) and [Scaleway Serverless Framework plugin](https://github.com/scaleway/serverless-scaleway-functions) to link and automate the deployment of some functions, get the associated endpoint and feed them in Kong configuration.

Pre-requisites:
- gnu make
- curl
- [jq](https://stedolan.github.io/jq/)
- [serverless cli](https://www.scaleway.com/en/docs/tutorials/install-serverless-framework/) (Serverless Compose is added as dependancy in `packages.json`)
- Docker Desktop to build and push the container.
- Having SCW_ACCESS_KEY and SCW_SECRET_KEY environment variables set.


serverless-compose.yml
----------------------

This file links all the `serverless.yml` services located in the different folders:

```
├── apiGateway
│   └── serverless.yml
├── getToken
│   └── serverless.yml
├── myApp
│   └── serverless.yml
└── serverless-compose.yml
```

It allows passing output from one service to an other as parameter.

myApp folder
------------

This folder contains two very simple Python functions to display html. The `commands` function display a link to `orders` using only `/orders` subpath.

Both these functions are deployed as "private". They can't be accessed without an `X-Auth-Token` header using a token which should be generated on the namesapce.


getToken folder
---------------

This folder contains a custom plugin to execute shell scripts using Serverless framework.

- The `deploy.sh` script will generate a token using Scaleway API and store it as a file. It uses `jq` to parse the JSON output.
- The `info.sh` script will output the token in a format expected by Serverless Compose.
- The `remove.sh` script only will delete the token file (the token is deleted together with the namespace when myApp service is deleted).


apiGateway folder
-----------------

This folder deploy a DB-less Kong container, build from the `Dockerfile` in the `kong` folder:

```
apiGateway
├── kong
│   ├── Dockerfile
│   ├── kong.conf
│   ├── kong.yml.template
│   └── start.sh
└── serverless.yml
```

[Kong](https://konghq.com/) is an API Gateway which can either be deployed with a database or with a static configuration file (DB-less mode).

The container is started with two environment variables (the HTTP endpoints of the functions) and a secret (the token). These are provided via Serverless Compose from the output of the `myApp` and `getToken` services.

The `kong.yml.template` is the DB-less configuration used by the API gateway and the `start.sh` script will parse it together with the environment variable to build the config file used by Kong.

The yaml file define routes (default as `/` and orders as `/orders`) and associate services to them:

```
routes:
- name: default-route
  paths:
  - /
  service: commands
- name: orders
  paths:
  - /orders
  service: orders
```

The services are specific url:

```
- name: commands
  url: https://${COMMANDS_URL}
- name: orders
  url: https://${ORDERS_URL}
```

The request-transformer plugin is used to automatically add the `X-Auth-Token` to all queries send to the services.

```
plugins:
- name: request-transformer
  config:
    add:
      headers:
      - X-Auth-Token:${TOKEN}
```

Other parameters can be added to the file to handle more use cases: https://docs.konghq.com/gateway/latest/production/deployment-topologies/db-less-and-declarative-config/

Running the example
-------------------

Just execute `make` to install the `packages.json` dependancies and execute the deployment.
