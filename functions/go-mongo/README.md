# Mongo + Go example

Example to demonstrate the possibility of connecting Serverless Functions to MongoDB.

For this example, [Scaleway Console](https://console.scaleway.com/) will be used for deployment.

> [!WARNING]
> This is a basic sample that does not use a certificate for authentication, not recommended for production.

## Requirements

- MongoDB created [documentation](https://www.scaleway.com/en/docs/managed-mongodb-databases/quickstart/#how-to-create-a-database-instance)
- The MongoDB instance must be public (not listening to a private network only since Serverless Functions do not support private networks yet)

## Step 1 - Mongo

After MongoDB creation, in the console, find the public endpoint of the database; it should look like: `<scw_database_id>.mgdb.<scw_region>.scw.cloud`

Once you get the endpoint, keep it somewhere for later use.

## Step 2 - Function Creation

Before creating the Function, we need to package it into a zip file.

> [!TIP]
> It's recommended to do it via command line tools, so the following example shows how to zip this current folder:
>
> ```sh
> cd ~/scaleway/serverless-examples/functions/go-mongo/
> zip -r go-mongo.zip *
> ```

- Create a Serverless Function namespace. [Documentation](https://www.scaleway.com/en/docs/serverless-functions/how-to/create-manage-delete-functions-namespace/#creating-a-serverless-functions-namespace)
- In the created namespace, create a Serverless Function.
- Select the latest Go runtime.
- Upload the previously created `go-mongo.zip`.
- Ensure the handler is `Handle`.
- Add required Secrets to the Function:

| Key                   | Value                                                                                                          |
| --------------------- | -------------------------------------------------------------------------------------------------------------- |
| MONGO_PUBLIC_ENDPOINT | (replace with required values): `<scw_database_id>.mgdb.<scw_region>.scw.cloud`                                  |
| MONGO_USER            | user created during MongoDB setup                                                                              |
| MONGO_PASSWORD        | password created during MongoDB user setup                                                                     |

## Step 3 - Test

Once your Serverless Function is in `ready` state, you can call it using the generated endpoint.

Result should be similar to:

```json
{ "ID": 2911126, "Name": "RandomName2911126" }
```

## Local testing

For testing, you can use [Go Offline Testing](https://github.com/scaleway/serverless-functions-go).

To run locally, execute: `go run cmd/main.go`.
