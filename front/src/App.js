import React from "react";
import { Router, Link } from "@reach/router";

import Home from "/views/Home/";
import Login from "/views/Login/";

import { isConnected } from "/utils/connect";

import "/App.scss";

const App = (
  <Router>
    <Home path="/" />
    <Login default={!isConnected()} path="/login" />
  </Router>
);

export default App;
