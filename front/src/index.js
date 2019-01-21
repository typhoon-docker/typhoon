import ReactDOM from "react-dom";
import { shouldMock } from "/utils/env";
import { importMocks } from "/utils/axios";

import App from "./App";

if (shouldMock) {
  importMocks().then(() => {
    ReactDOM.render(App, document.getElementById("app"));
  });
} else {
  ReactDOM.render(App, document.getElementById("app"));
}
