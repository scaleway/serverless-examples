// import type { AWS } from "@serverless/typescript";

import functions from "@functions/index";
import ScalewayServerlessConfiguration from "@libs/scalewayServerless";

const serverlessConfiguration: ScalewayServerlessConfiguration = {
  service: "nodejs-typescript-test",
  frameworkVersion: "3",
  configValidationMode: "off",
  plugins: ["serverless-esbuild", "serverless-scaleway-functions"],

  provider: {
    name: "scaleway",
    runtime: "node18",
    swcRegion: "fr-par",
    env: {
      AWS_NODEJS_CONNECTION_REUSE_ENABLED: "1",
      NODE_OPTIONS: "--enable-source-maps --stack-trace-limit=1000",
    },
  },
  // import the function via Functions index.ts file
  functions: functions,

  custom: {
    esbuild: {
      keepOutputDirectory: true,
      bundle: true,
      minify: false,
      sourcemap: true,
      exclude: ["aws-sdk"],
      format: "esm",
      define: { "require.resolve": undefined },
      platform: "node",
      concurrency: 10,
    },
  },
};

// Log to check output of serverless.ts
// console.log(JSON.stringify(serverlessConfiguration, null, 2));

module.exports = serverlessConfiguration;
