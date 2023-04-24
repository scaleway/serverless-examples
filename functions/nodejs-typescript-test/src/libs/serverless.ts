type SwcRegion = "fr-par" | "nl-ams" | "pl-waw";

type Runtime =
  | "node19"
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

type MemoryLimit = "128" | "256" | "512" | "1024" | "2048" | "3072" | "4096";

type ConfigValidationMode = "warn" | "error" | "off";

type Events = any;

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

export interface ScalewayFunction {
  handler: string;
  env?: Record<string, string>;
  secret?: Record<string, string>;
  minScale?: number;
  maxScale?: number;
  maxConcurrency?: number;
  memoryLimit?: MemoryLimit;
  timeout?: number;
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
