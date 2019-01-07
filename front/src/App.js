import React from "react";
import { Router, Link } from "@reach/router";

import Login from "/views/Login/";

import "/App.scss";

const App = (
  <Router>
    <Login path="/" />
  </Router>
);

export default App;
