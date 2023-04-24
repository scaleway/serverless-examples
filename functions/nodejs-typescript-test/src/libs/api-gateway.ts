import type {
  APIGatewayProxyEvent,
  APIGatewayProxyResult,
  Handler,
} from "aws-lambda";
import type { FromSchema } from "json-schema-to-ts";

type ValidatedAPIGatewayProxyEvent<S> = Omit<APIGatewayProxyEvent, "body"> & {
  body: FromSchema<S>;
};
export type ValidatedEventAPIGatewayProxyEvent<S> = Handler<
  ValidatedAPIGatewayProxyEvent<S>,
  APIGatewayProxyResult
>;

interface IformatJSONResponse {
  statusCode: number;
  response: Record<string, unknown>;
}

const thing = () => {
  return "nothing to see here";
};

export const formatJSONResponse = (props: IformatJSONResponse) => {
  return {
    statusCode: props.statusCode || 200,
    body: JSON.stringify(props.response),
  };
};
