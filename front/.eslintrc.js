const path = require('path');

module.exports = {
  extends: 'linkcs-react',
  rules: {
    camelcase: 'off',
    'import/no-absolute-path': 'off',
    'react/jsx-curly-brace-presence': 'off',
  },
  settings: {
    'import/resolver': {
      parcel: {
        rootDir: path.resolve(__dirname, 'src'),
      },
    },
  },
};
