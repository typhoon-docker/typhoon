import React, { Fragment } from 'react';
import { Router } from '@reach/router';
import 'highlight.js/styles/atom-one-dark.css';

import Home from '/views/Home';
import Admin from '/views/Admin';
import Login from '/views/Login/';
import New from '/views/New/';
import Project from '/views/Project/';
import CallbackViaRezo from '/views/CallbackViaRezo/';
import CallbackGitHub from '/views/CallbackGitHub/';
import AOA from '/views/404/';

import { useIsConnected } from '/utils/connect';

import '/App.css';

const App = () => {
  const isConnected = useIsConnected();
  return (
    <Router>
      <CallbackViaRezo path="/callback/viarezo" />
      {isConnected ? (
        <Fragment default>
          <Home path="/" />
          <Admin path="/admin" />
          <New path="/new" />
          <Project path="/project/:projectID" />
          <CallbackGitHub path="/callback/github" />
          <AOA default />
        </Fragment>
      ) : (
        <Login default />
      )}
    </Router>
  );
};

export default App;
