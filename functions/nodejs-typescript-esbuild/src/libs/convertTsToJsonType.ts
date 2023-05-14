import * as TJS from "typescript-json-schema";

import { resolve } from "path";

interface IgenerateJsonSchema {
  path: string;
  typeName: string;
}

export default function generateJSONSchema(props: IgenerateJsonSchema) {
  // optionally pass argument to schema generator
  const settings: TJS.PartialArgs = {
    required: true,
  };
  const { path, typeName } = props;
  const basePath = "../../";
  // optionally pass ts compiler options
  // const compilerOptions: TJS.CompilerOptions = {
  //   strictNullChecks: true,
  // };

  // optionally pass a base path

  const program = TJS.getProgramFromFiles(
    [resolve(`${path}/schema.ts`)],
    // compilerOptions,
    basePath
  );

  // get the schema and return it
  const schema = TJS.generateSchema(program, typeName, settings);

  return schema;
}
