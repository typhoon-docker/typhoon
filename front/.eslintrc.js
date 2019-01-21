const path = require("path");

module.exports = {
  extends: "linkcs-react",
  rules: {
    camelcase: "off",
    "import/no-absolute-path": "off"
  },
  settings: {
    "import/resolver": {
      parcel: {
        rootDir: path.resolve(__dirname, "src")
      }
    }
  }
};
