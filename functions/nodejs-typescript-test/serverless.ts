import type { AWS } from "@serverless/typescript";

import hello from "@functions/hello";

const serverlessConfiguration: AWS = {
  service: "nodejs-typescript-test",
  frameworkVersion: "3",
  plugins: ["serverless-esbuild"],
  package: { individually: true },
  provider: {
    name: "aws",
    runtime: "nodejs14.x",
    apiGateway: {
      minimumCompressionSize: 1024,
      shouldStartNameWithService: true,
    },
    environment: {
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
