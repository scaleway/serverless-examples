// import type { AWS } from "@serverless/typescript";

import hello from "@functions/hello";
import ScalewayServerlessConfiguration from "@libs/serverless";

const serverlessConfiguration: ScalewayServerlessConfiguration = {
  service: "nodejs-typescript-test",
  frameworkVersion: "3",
  configValidationMode: "off",
  plugins: ["serverless-esbuild", "serverless-scaleway-functions"],
  package: { individually: false },
  provider: {
    name: "scaleway",
    runtime: "node14",
    swcRegion: "fr-par",
    // apiGateway: {
    //   minimumCompressionSize: 1024,
    //   shouldStartNameWithService: true,
    // },
    env: {
      AWS_NODEJS_CONNECTION_REUSE_ENABLED: "1",
      NODE_OPTIONS: "--enable-source-maps --stack-trace-limit=1000",
    },
  },
  // import the function via paths
  functions: { hello },

  custom: {
    esbuild: {
      keepOutputDirectory: true,
      bundle: true,
      minify: false,
      sourcemap: true,
      exclude: ["aws-sdk"],
      target: "node18",
      define: { "require.resolve": undefined },
      platform: "node",
      concurrency: 10,
    },
  },
};

module.exports = serverlessConfiguration;
