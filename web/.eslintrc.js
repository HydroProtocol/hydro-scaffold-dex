module.exports = {
  env: {
    browser: true,
    node: true,
    es6: true
  },
  extends: ['eslint:recommended', 'plugin:react/recommended'],
  globals: {
    Atomics: 'readonly',
    SharedArrayBuffer: 'readonly'
  },
  parser: 'babel-eslint',
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
      experimentalObjectRestSpread: true
    },
    ecmaVersion: 2018,
    sourceType: 'module'
  },
  plugins: ['react'],
  rules: {
    'no-console': 'warn',
    'no-unused-vars': 'warn',
    'react/prop-types': 'off',
    'no-case-declarations': 'off',
    'no-extra-boolean-cast': 'off',
    'react/no-unescaped-entities': 'off'
  }
};
