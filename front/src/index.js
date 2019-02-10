import React from 'react';
import ReactDOM from 'react-dom';
import { shouldMock } from '/utils/env';
import { importMocks as importMocksTyphoon } from '/utils/typhoonAPI';
import { importMocks as importMocksGitHub } from '/utils/githubAPI';

import App from './App';

if (shouldMock) {
  Promise.all([importMocksTyphoon(), importMocksGitHub()]).then(() => {
    ReactDOM.render(<App />, document.getElementById('app'));
  });
} else {
  ReactDOM.render(<App />, document.getElementById('app'));
}
