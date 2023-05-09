type SwcRegion = "fr-par" | "nl-ams" | "pl-waw";

type Runtime =
  | "node20"
  | "node18"
  | "node17"
  | "node16"
  | "node14"
  | "node10"
  | "node8"
  | "python311"
  | "python310"
  | "python39"
  | "python38"
  | "python37"
  | "python"
  | "python3"
  | "php82"
  | "go120"
  | "go119"
  | "go118"
  | "go117"
  | "go113"
  | "golang"
  | "rust165";

type MemoryLimit = 128 | 256 | 512 | 1024 | 2048 | 3072 | 4096;

type ConfigValidationMode = "warn" | "error" | "off";

type Events = any; // needs to be typed!

interface Package {
  individually?: false; // Not supported by scaleway
  patterns?: Array<string>;
}

interface ScalewayServerlessProvider {
  name: "scaleway";
  runtime: Runtime;
  env?: Record<string, string>;
  secret?: Record<string, string>;

  swcRegion: SwcRegion;
}

interface ScalewayFunctionTimeOut {
  seconds: number;
}
type ScalewayFunctionScales =
  | 0
  | 1
  | 2
  | 3
  | 4
  | 5
  | 6
  | 7
  | 8
  | 9
  | 10
  | 11
  | 12
  | 13
  | 14
  | 15
  | 16
  | 17
  | 18
  | 19
  | 20;
export interface ScalewayFunction {
  handler: string;
  description?: string;
  env?: Record<string, string>;
  secret?: Record<string, string>;
  minScale?: ScalewayFunctionScales;
  maxScale?: ScalewayFunctionScales;
  maxConcurrency?: number;
  memoryLimit?: MemoryLimit;
  // timeout?: any; NOT WORKING! 
  runtime?: Runtime;
  events?: Events;
  httpOption?: "enabled" | "redirected"; // Force HTTPS redirection
  privacy?: "public" | "private"; // (Optional): defines whether a function may be executed anonymously (public) or only via an authentication mechanism (private) (default: public)
}

interface ScalewayFunctionMap {
  [key: string]: ScalewayFunction;
}

export default interface ScalewayServerlessConfiguration {
  service: string;
  frameworkVersion: string;
  configValidationMode?: ConfigValidationMode;
  useDotenv?: boolean;
  provider: ScalewayServerlessProvider;
  plugins: Array<string>;
  package?: Package;
  custom?: Record<string, any>;
  functions: ScalewayFunctionMap;
}

interface HandlerResponse {
  statusCode?: Number;
  headers?: Record<string, unknown>;
  body: Record<string, unknown>;
}

export const formatScalewayHandlerJSONResponse = (
  response: HandlerResponse
) => {
  return {
    headers: {
      "Content-Type": "application/json",
    },
    statusCode: response.statusCode ? response.statusCode : 200,
    body: JSON.stringify(response.body),
  };
};

export interface AllScalewayFunctionData {
  [functionName: string]: ScalewayFunction;
}

// Type definitions for non-npm package scaleway-functions 1.0
// Project: https://www.scaleway.com/en/serverless-functions/
// Definitions by: MrMicky <https://github.com/MrMicky-FR>
// Definitions: https://github.com/DefinitelyTyped/DefinitelyTyped

export type Handler<TResult = Response | object> = (
  event: Event,
  context: Context,
  callback: Callback<TResult>
) => void | TResult | Promise<TResult>;

export type Callback<TResult = Response | object> = (
  error?: Error | string | null,
  result?: TResult
) => void;

// https://github.com/scaleway/scaleway-functions-runtimes/blob/master/events/context.go
export interface Context {
  memoryLimitInMb: number;
  functionName: string;
  functionVersion: string;
}

// https://github.com/scaleway/scaleway-functions-runtimes/blob/master/events/http.go
export interface Event {
  path: string;
  httpMethod: string;
  headers: Record<string, string>;
  queryStringParameters: Record<string, string>;
  stageVariables: Record<string, string>;
  body: unknown;
  isBase64Encoded: boolean;
  requestContext: RequestContext;
}

export interface RequestContext {
  stage: string;
  httpMethod: string;
}

// https://github.com/scaleway/scaleway-functions-runtimes/blob/master/handler/utils.go
export interface Response {
  statusCode: number;
  body?: string | object;
  headers?: Record<string, string>;
  isBase64Encoded?: boolean;
}
