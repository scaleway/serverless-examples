import { middyfy } from "@libs/middyfy";
import {
  Handler,
  ValidatedScalewayHandlerEvent,
  formatScalewayHandlerJSONResponse,
} from "@libs/scalewayServerless";

import { helloWorld } from "hello-world-npm";
import { Response } from "scaleway-functions";

import { ISchema } from "./schema";

export const hello: Handler = async (
  event: ValidatedScalewayHandlerEvent<ISchema>
): Promise<Response> => {
  const data = event?.body;
  const timeToWait = 500;
  await new Promise((res) => setTimeout(res, timeToWait));

  return formatScalewayHandlerJSONResponse({
    statusCode: 200,
    headers: {
      "Content-Type": "application/json",
    },
    body: {
      message: {
        one: `Hello ${
          data?.name || "person"
        }, welcome to the exciting Serverless world!`,
        two: helloWorld(),
        env: process.env.TEST_2,
      },
    },
  });
};
export const handler = middyfy(hello); // Its important that the exported funtion name is handler.
