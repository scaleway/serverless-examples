module.exports = {
  root: true,
  extends: ["@serverless/eslint-config"],
  rules: {
    "import/prefer-default-export": "off",
  },

  ignorePatterns: ["node_modules", ".turbo", ".next"],
};
