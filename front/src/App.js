import React, { Fragment } from 'react';
import { Router } from '@reach/router';

import Home from '/views/Home';
import Admin from '/views/Admin';
import Login from '/views/Login/';
import Callback from '/views/Callback/';
import AOA from '/views/404/';

import { useIsConnected } from '/utils/connect';

import '/App.css';

const App = () => {
  const isConnected = useIsConnected();
  return (
    <Router>
      <Callback path="/callback/viarezo" />
      {isConnected ? (
        <Fragment default>
          <Home path="/" />
          <Admin path="/admin" />
          <AOA default />
        </Fragment>
      ) : (
        <Login default />
      )}
    </Router>
  );
};

export default App;
