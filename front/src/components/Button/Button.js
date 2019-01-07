import React from "react";

import { button } from "./Button.css";

const Button = ({ color, className: cn, style: s, ...props }) => {
  const style = s ? { ...s } : {};
  style.backgroundColor = `var(--${color || "text"})`;

  const className = cn ? [cn, button] : button;
  return <button {...props} className={className} style={style} />;
};

export default Button;
