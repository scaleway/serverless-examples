module.exports = {
  root: true,
  // extends: ['@tif/eslint-config/serverless'],
  rules: {
    "import/prefer-default-export": "off",
  },

  ignorePatterns: ["node_modules", ".turbo", ".next"],
};
