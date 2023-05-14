/**
 * Type incoming data here. This data will be used to type the body inside of the function.
 * It will also be converted to JSON typings in order to be able to use it to
 * validate incoming data in an API gateway.
 */

import generateJSONSchema from "../../libs/convertTsToJsonType";
import { handlerPath } from "../../libs/handler-resolver";

// Write out shape of incomming dat in TS Interface
export interface ISchema {
  name?: string;
}

// Convert to JSON typings for API gateway checking
const path = `${handlerPath(__dirname)}`;

const schema = generateJSONSchema({ path, typeName: "ISchema" });

// console.log(JSON.stringify(schema, null, 2));

export default schema;
