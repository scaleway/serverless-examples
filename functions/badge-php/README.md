# Repository Badge PHP example

This example demonstrates how PHP functions work using a dependency from composer.
With this example, you can generate a simple badge with the text, status and color of your choice, just like this
one: ![PHP Example: OK](https://serverlessexamplesmpjacsh6-badge.functions.fnc.fr-par.scw.cloud/?text=PHP%20Example&status=OK&color=purple)

## Requirements

This example assumes you are familiar with how serverless functions work. If needed, you can
check [Scaleway official documentation](https://www.scaleway.com/en/docs/serverless/functions/quickstart/)

This example uses the Scaleway Serverless Framework Plugin. Please set up your environment with the requirements stated
in the [Scaleway Serverless Framework Plugin](https://github.com/scaleway/serverless-scaleway-functions) before trying
out the example.

## In this example

This example shows to create a basic PHP function with dependencies provided in a composer.json file.
The function generates an SVG image using Twig based on user input.

## Setup

Once your environment is set up, you can run:

```console
npm install

serverless deploy
```

## Running

Then, from the given URL, you can run open your browser and navigate to it.
You should see the default
example: ![Scaleway: approved](https://serverlessexamplesmpjacsh6-badge.functions.fnc.fr-par.scw.cloud/)

You can tweak the following query parameters to make this badge your own:

| key    | description                                                                                                                   |
|--------|-------------------------------------------------------------------------------------------------------------------------------|
| text   | The text you want the badge to show. Default: Scaleway                                                                        |
| status | The status text you want the badge to show. Default: approved                                                                 |
| color  | The color of the status part. Default: green. **This exemple does not accept hex color code, please use their literal names** |
