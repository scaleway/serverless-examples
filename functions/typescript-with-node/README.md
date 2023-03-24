# Use typescript with Node runtime

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway's official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

Finally, you will need Node.js installed in your computer to run this example.

## Context

By default, Node runtime treats files with the .js suffix. If you wish to use Typescript language with Node runtime, you can do so by following this example.

## Description

This example aims to show how to use Typescript language with Node runtime (node 18 runtime in this example). Used packages are specified in `package.json`.

The function in this example returns a simple "Hello world!" with a status code 200.

## Setup

### Install npm modules

Once your environment is set up (see [Requirements](#requirements)), you can install `npm` dependencies from `package.json` file using:

```shell
npm install
```

### Install a Typescript compiler

First, it is necessary to install the [Typescript compiler package](https://www.npmjs.com/package/typescript) globally.

```shell
npm install -g typescript
```

You can run `tsc --version` to ensure the compiler is correctly installed.

### Create a Typescript configuration file

When this is done, you can initialize the Typescript project with Node.js. For that, you can run:

```shell
tsc --init
```

This will create a `tsconfig.json` file in the project root directory.

### Transpile your code

Before deploying your function, you need to transpile your Typescript code into brower readable JavaScript.

```shell
tsc
```

### Deploy your function

Finally, you can deploy your function with:

```shell
serverless deploy
```

## Running

Once your function is deployed, you can check the result in a browser or by running the following command:

```console
# Get request
curl -i -X GET <function URL>
```

The output should be "Hello world!".
