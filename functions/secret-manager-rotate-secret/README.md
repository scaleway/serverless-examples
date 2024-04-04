# Rotate RDB Credentials

This function will rotate the credentials of and RDB database, stored in a secret in the Secret Manager.

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

An RDB database is required for this to work. The credentials of this database MUST be stored in a secret with the following layout:
```json
{
    "engine": "postgres|mysql",
    "username": "db_username",
    "password": "db_password",
    "host": "db_ip_or_hostname",
    "dbname": "db_name",
    "port": "db_port"
}
```

For the function to work, it needs an API key with the following permissions:
- `SecretManagerFullAccess`
- `RelationalDatabasesFullAccess`

## Context

The function will generate a new password for your RDB credentials using the Scaleway API, then it will access the secret where it is stored. After that it will update the RDB credentials with the `username` configured in the secret and the new generated password. Finally it will create a new version of the secret with the new credentials.

## Setup

You will need to adjust to your needs the `env` and `secret` settings in the `serverless.yml` file.

Once your environment is set up, you can run:

```console
npm install

serverless deploy
```

## Running

You can use `curl` to trigger your function. It requires the following input, you can store it in a file called `req.json`.
```json
{
    "rdb_instance_id": "your RDB instance ID",
    "secret_id": "the secret ID where credentials are stored"
}
```

```console
curl <function URL> -d @req.json
```

**Update with the expected output**
The result should be similar to:

```console
HTTP/2 200
content-length: 21
content-type: text/plain
date: Tue, 17 Jan 2023 14:02:46 GMT
server: envoy
x-envoy-upstream-service-time: 222

database credentials updated%
```
