# Interact with a PostgreSQL database using functions

This example shows how to connect to a managed PostgreSQL database and perform a query on it.

## Requirements

This example assumes that you are familiar with some products of Scaleway's ecosystem:

* how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/).
* how a managed database for PostgreSQL works, and especially how to create a database and create users with appropriate permissions. Please refer to scaleway's documentation [here](https://www.scaleway.com/en/docs/managed-databases/postgresql-and-mysql/quickstart/).

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Additionnaly it uses the [serverless-functions-node](https://github.com/scaleway/serverless-functions-node) library for local testing.

## Context

This example shows how to connect to a managed PostgreSQL database and perform a query on it. This example can be extended to adding, deleting, or processing data within a database.

## Description

The function connects to a PostgreSQL database and performs an example query on it. This example uses Node 18 runtime. Used packages are specified in `package.json`.

## Setup

### Create a managed PostgreSQL database

Create a PostgreSQL database and a user profile with appropriate access permissions.

### Fill environment variables

Fill your secrets within `serverless.yml` file:

```yml
secret:
    PG_HOST: "your host IP address"
    PG_USER: "your database username"
    PG_DATABASE: "your database name"
    PG_PASSWORD: "your databse user password"
    PG_PORT: "your database port"
```

### Install npm modules

Once your environment is set up, you can install `npm` dependencies from `package.json` file using:

```sh
npm install
```

### Test locally

Once your environment is set up, you can run:

```sh
NODE_ENV=test node handler.js
```

This will launch a local server, allowing you to test the function. Then, you can run in another terminal:

```sh
curl -X GET http://localhost:8080
```

The output should be similar to:

```sh
[{"user":"<PG_USER>"}]
```

## Deploy and run

Finally, if the test succeeded, you can deploy your function with:

```sh
serverless deploy
```

Then, from the given URL, you can run:

```sh
curl -X GET <function URL>
```

When invoking this function, the output should be similar to the one obtained when testing locally.
