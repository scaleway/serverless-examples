# Serverless - Scaleway Node.js Typescript ESbuild

This project has been generated using the `aws-nodejs-typescript` template from the [Serverless framework](https://www.serverless.com/), and has then been modified to run and build for Scaleway serverless functions.

For detailed instructions, please refer to the [documentation](https://www.scaleway.com/en/docs/faq/serverless-functions/).

## Installation/deployment instructions

Depending on your preferred package manager, follow the instructions below to deploy your project.

> **Requirements**: NodeJS `(v.16)+`. If you're using [nvm](https://github.com/nvm-sh/nvm), run `nvm use` to ensure you're using the same Node version in local and in your lambda's runtime.

### Using NPM

- Run `npm i` to install the project dependencies.
- Run `npx sls deploy` to deploy this stack to Scaleway Functions.

### Using Yarn

- Run `yarn` to install the project dependencies.
- Run `yarn sls deploy` to deploy this stack to Scaleway Functions.

## Test your service

This template contains a single lambda function triggered by an HTTP request made on the provisioned serverless function.

Sending a `POST` request to the URL with a payload containing a string property named `name` will result in API Gateway returning a `200` HTTP status code with a message saluting the provided name and the detailed event processed by the lambda.

> :warning: As is, this template, once deployed, opens a **public** endpoint within your Scalway account resources. Anybody with the URL can actively execute the API Gateway endpoint and the corresponding lambda. You should protect this endpoint with the authentication method of your choice. An example can be found in the hello-2 handler.

### Locally

Local testing is not yet functional.

In order to test the hello function locally, run the following command:

- `npx sls invoke local -f hello --path src/functions/hello/mock.json` if you're using NPM.
- `yarn sls invoke local -f hello --path src/functions/hello/mock.json` if you're using Yarn.

Check the [sls invoke local command documentation](https://www.serverless.com/framework/docs/providers/aws/cli-reference/invoke-local/) for more information.

### Remotely

Copy and replace your `url`, found in Serverless `deploy` command output, and `name` parameter in the following `curl` command in your terminal or in Postman to test your newly deployed application.

```
curl --location --request POST 'https://myApiEndpoint/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Frederic"
}'
```

## Deployment instructions

To deploy the template in its current form, you need to either set the environment variables or add the project and key to the CLI parameters. Due to a bug in the Scaleway plugin, the credentials inside of the serverless.ts file won't work (yet).

## Template features

- ESbuild bundling minification
- Typescript with types
- Dynamic serverless.ts config

### Project structure

The project code base is mainly located within the `src` folder. This folder is divided into:

- `functions` - containing code base and configuration for your lambda functions.
- `libs` - containing shared code base between your lambdas.

```
.
├── src
│   ├── functions                   # Lambda configuration and source code folder
│   │   ├── hello                   # `hello` example
│   │   │   ├── handler.ts          # `Hello` lambda source code
│   │   │   ├── index.ts            # `Hello` lambda Serverless configuration
│   │   │   ├── mock.json           # `Hello` lambda input parameter, if any, for local invocation
│   │   │   └── schema.ts           # Typescript schema for incoming data (used for typechecking and later on for API gateway)
│   │   ├── hello-2                 # `hello-2` example with external deps and async await and privacy enabled
│   │   │   ├── handler.ts          # `Hello` lambda source code
│   │   │   ├── index.ts            # `Hello` lambda Serverless configuration
│   │   │   ├── mock.json           # `Hello` lambda input parameter, if any, for local invocation
│   │   │   └── schema.ts           # Typescript schema for incoming data (used for typechecking and later on for API gateway)
│   │   └── index.ts                # Import/export of all lambda configurations (looking for a way to automate this)
│   └── libs                        # Lambda shared code
│       └── convertTsToJsonType     # API Gateway specific helpers
│       └── handlerResolver.ts      # Sharable library for resolving lambda handlers directories
│       └── scalewayServerless.ts   # Scaleway's typings and helpers
│       └── scalewayServerless.ts   # Scaleway's typings and helpers
│
│
├── package.json
├── .eslintrc.js                # Eslint config
├── serverless.ts               # Serverless service file
├── tsconfig.json               # Typescript compiler configuration
├── tsconfig.paths.json         # Typescript paths
```

# Future Features

- [ ] Auto import the functions in /src/functions/index.ts
- [ ] support python
- [ ] support golang
