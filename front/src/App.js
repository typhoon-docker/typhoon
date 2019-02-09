import React from 'react';
import { Router } from '@reach/router';

import Home from '/views/Home/';
import Login from '/views/Login/';
import Callback from '/views/Callback/';

import { isConnected } from '/utils/connect';

import '/App.css';

const App = (
  <Router>
    <Home path="/" />
    <Login default={!isConnected()} path="/login" />
    <Callback path="/callback/viarezo" />
  </Router>
);

export default App;
