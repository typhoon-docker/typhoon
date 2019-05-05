import h from '/utils/h';

import ReactDOM from 'react-dom';
import { shouldMock } from '/utils/env';

import App from './App';

if (shouldMock) {
  Promise.all([
    import('/utils/typhoonAPI').then(({ importMocks }) => importMocks()),
    import('/utils/githubAPI').then(({ importMocks }) => importMocks()),
  ]).then(() => {
    ReactDOM.render(<App />, document.getElementById('app'));
  });
} else {
  ReactDOM.render(<App />, document.getElementById('app'));
}
