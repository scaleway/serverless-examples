# Python chatbot example

This example shows how to handle different HTTP methods (POST and GET) as well as specific dependencies.

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying out the example.

## Example explanation

**Context:** This example shows how to handle GET and POST methods. For this example, [Chatterbot](https://github.com/gunthercox/ChatterBot/tree/1.0.4) (a machine learning, conversational dialog engine) was used. The model used here has already been trained to communicate in english to speed up the installation. The results of the training are stored in a sqlite database (in `app/english-corpus.sqlite3`). If you wish to train your own model, you can follow the instructions given in [Chatterbot documentation](https://chatterbot.readthedocs.io/en/stable/training.html).

**Explanation:** When the function is triggered by a GET method, it renders an HTML file. When the function is triggered by a POST method with the parameter "message", a response is given by the trained chatbot. In all other cases (method not handled or parameter missing), an error is thrown.

## Setup

Once your environment is set up, you can run:

```console
# Install node dependencies
npm install

# Deploy
./bin/deploy.sh
```

## Running

Then, you can test your function by sending the following request:

```console
# POST request
curl -i -X POST <function URL> -d '{"message":"Hello"}'
```

This will tell the chatbot "Hello". The expected answer should be something similar to "Hello" or "Hi".

The result of your function can also be checked through a browser. The expected behavior is to render the HTML file. And whenever a message is sent to the chatbot, a response should appear in the conversation.
