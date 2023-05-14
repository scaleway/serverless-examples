import {
  Handler,
  ValidatedScalewayHandlerEvent,
  formatScalewayHandlerJSONResponse,
} from "@libs/scalewayServerless";
import { ISchema } from "./schema";
import { middyfy } from "@libs/middyfy";

const hello: Handler = async (
  event: ValidatedScalewayHandlerEvent<ISchema>
) => {
  const data = event?.body;

  return formatScalewayHandlerJSONResponse({
    statusCode: 200,
    body: {
      message: `Hello ${
        data.name || "person"
      }, welcome to the exciting Serverless world!`,
      event: data,
    },
  });
};

export const handler = middyfy(hello);
