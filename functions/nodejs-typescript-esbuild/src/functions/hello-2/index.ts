// import schema from './schema';
import { handlerPath } from "@libs/handler-resolver";
import { ScalewayFunction } from "@libs/scalewayServerless";
// import schema from "./schema";

const hello2: ScalewayFunction = {
  handler: `${handlerPath(__dirname)}/handler.handler`,
  memoryLimit: 128,
  privacy: "private",
  description: "test description",
  httpOption: "redirected",
  minScale: null,
  maxScale: null,
  env: { TEST: "test", TEST_2: "test2" },

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
// console.log(JSON.stringify(schema, null, 2));

export default hello2;
