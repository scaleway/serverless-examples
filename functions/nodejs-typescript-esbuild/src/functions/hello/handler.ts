import {
  Handler,
  formatScalewayHandlerJSONResponse,
} from "@libs/scalewayServerless";

export const handler: Handler = async (event) => {
  const data = event.body as Record<string, unknown>;

  return formatScalewayHandlerJSONResponse({
    statusCode: 200,
    body: {
      message: `Hello ${
        data?.name || "person"
      }, welcome to the exciting Serverless world!`,
    },
  });
};
