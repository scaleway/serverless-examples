// import schema from './schema';
import { handlerPath } from "@libs/handler-resolver";
import { ScalewayFunction } from "@libs/scalewayServerless";

const hello: ScalewayFunction = {
  handler: `${handlerPath(__dirname)}/handler.handler`,
  memoryLimit: 128,
  events: [
    // {
    //   http: {
    //  HTTP events to deploy a API gateway are not yet supported by scaleway here for illustration only
    //     method: 'post',
    //     path: 'hello',
    //     request: {
    //       schemas: {
    //         'application/json': schema,
    //       },
    //     },
    //   },
    // },
  ],
};

export default hello;
