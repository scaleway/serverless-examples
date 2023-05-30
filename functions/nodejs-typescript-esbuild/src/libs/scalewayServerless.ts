type SwcRegion = "fr-par" | "nl-ams" | "pl-waw"

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
  | "rust165"

type MemoryLimit = 128 | 256 | 512 | 1024 | 2048 | 3072 | 4096

type ConfigValidationMode = "warn" | "error" | "off"

type Events = any // needs to be typed!

interface Package {
  individually?: false // Not supported by scaleway
  patterns?: Array<string>
}

interface ScalewayServerlessProvider {
  name: "scaleway"
  runtime: Runtime
  env?: Record<string, string>
  secret?: Record<string, string>

  swcRegion: SwcRegion
}

type ScalewayFunctionTimeOut = string

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
  | 20
export interface ScalewayFunction {
  handler: string
  description?: string
  env?: Record<string, string>
  secret?: Record<string, string>
  minScale?: ScalewayFunctionScales
  maxScale?: ScalewayFunctionScales
  maxConcurrency?: number
  memoryLimit?: MemoryLimit
  timeout?: ScalewayFunctionTimeOut
  runtime?: Runtime
  events?: Events
  httpOption?: "enabled" | "redirected" // Force HTTPS redirection
  privacy?: "public" | "private" // (Optional): defines whether a function may be executed anonymously (public) or only via an authentication mechanism (private) (default: public)
}

interface ScalewayFunctionMap {
  [key: string]: ScalewayFunction
}

export default interface ScalewayServerlessConfiguration {
  service: string
  frameworkVersion: string
  configValidationMode?: ConfigValidationMode
  useDotenv?: boolean
  provider: ScalewayServerlessProvider
  plugins: Array<string>
  package?: Package
  custom?: Record<string, any>
  functions: ScalewayFunctionMap
}

export const formatScalewayHandlerJSONResponse = (
  response: ScalewayHandlerResponse
) => {
  return {
    headers: response.headers,
    statusCode: response.statusCode,
    body: JSON.stringify(response.body),
  }
}

export interface AllScalewayFunctionData {
  [functionName: string]: ScalewayFunction
}

// Type definitions for non-npm package scaleway-functions 1.0
// Project: https://www.scaleway.com/en/serverless-functions/
// Definitions by: MrMicky <https://github.com/MrMicky-FR>
// Definitions: https://github.com/DefinitelyTyped/DefinitelyTyped

/**
 * @example <caption>Defining a custom handler type</caption>
 * import { Handler } from ''
 *
 * interface NameEvent {
 *     fullName: string
 * }
 * interface NameResult {
 *     firstName: string
 *     middleNames: string
 *     lastName: string
 * }
 * type PersonHandler = Handler<NameEvent, NameResult>
 *
 * export const handler: PersonHandler = async (event) => {
 *   const names = event.fullName.split(' ')
 *   const firstName = names.shift()
 *   const lastName = names.pop()
 *   return { firstName, middleNames: names, lastName }
 * }
 *
 *
 * @param event
 *      Parsed JSON data in the lambda request payload. For an ScaleWay service triggered
 *      lambda this should be in the format of a type ending in Event, for example the
 *      S3Handler receives an event of type S3Event.
 * @param context
 *      Runtime contextual information of the current invocation, for example the caller
 *      identity, available memory and time remaining, legacy completion callbacks, and
 *      a mutable property controlling when the lambda execution completes.
 * @param callback
 *      NodeJS-style completion callback that the Scaleway Lambda runtime will provide that can
 *      be used to provide the lambda result payload value, or any execution error. Can
 *      instead return a promise that resolves with the result payload value or rejects
 *      with the execution error.
 * @return
 *      A promise that resolves with the lambda result payload value, or rejects with the
 *      execution error. Note that if you implement your handler as an async function,
 *      you will automatically return a promise that will resolve with a returned value,
 *      or reject with a thrown value.
 */
export type Handler<TEvent = any, TResult = any> = (
  event: TEvent,
  context: Context,
  callback: Callback<TResult>
) => void | Promise<TResult>

export type Callback<TResult = ScalewayHandlerResponse | object> = (
  error?: Error | string | null,
  result?: TResult
) => void

// https://github.com/scaleway/scaleway-functions-runtimes/blob/master/events/context.go
export interface Context {
  memoryLimitInMb: number
  functionName: string
  functionVersion: string
}

// https://github.com/scaleway/scaleway-functions-runtimes/blob/master/events/http.go
export interface ScalewayHandlerEvent {
  path: string
  httpMethod: string
  headers: Record<string, string>
  queryStringParameters: Record<string, string>
  stageVariables: Record<string, string>
  body: Record<string, unknown>
  isBase64Encoded: boolean
  requestContext: RequestContext
}

export type ValidatedScalewayHandlerEvent<S extends object> = Omit<
  ScalewayHandlerEvent,
  "body"
> & { body: S }

export interface RequestContext {
  stage: string
  httpMethod: string
}

// https://github.com/scaleway/scaleway-functions-runtimes/blob/master/handler/utils.go
export interface ScalewayHandlerResponse {
  statusCode?: number
  body?: object
  headers?: Record<string, string>
  isBase64Encoded?: boolean
}
