import React from "react";
import { Router, Link } from "@reach/router";

import Login from "/views/Login/";

import { isConnected } from "/utils/connect";

import "/App.scss";

const App = (
  <Router>
    <Login default={!isConnected()} path="/login" />
  </Router>
);

export default App;
