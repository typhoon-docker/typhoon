import React from "react";

import { button } from "./Button.css";

const Button = ({ className: cn, ...props }) => {
  const className = cn ? [cn, button] : button;
  return <button {...props} className={className} />;
};

export default Button;
