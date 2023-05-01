import {
  Handler,
  formatScalewayHandlerJSONResponse,
} from "@libs/scalewayServerless";

import { helloWorld } from "hello-world-npm";

export const handler: Handler = async (event) => {
  const data = event.body as Record<string, unknown>;

  return formatScalewayHandlerJSONResponse({
    statusCode: 200,
    body: {
      message: {
        one: `Hello ${data.name}, welcome to the exciting Serverless world!`,
        two: helloWorld(),
      },
    },
    //event,
  });
};
