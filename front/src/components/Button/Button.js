import React from 'react';

import { button } from './Button.css';

const Button = ({ color, className: cn, style: s, ...props }) => {
  const style = s ? { ...s } : {};
  style.backgroundColor = `rgb(var(--${color || 'text'}))`;

  const className = [button];
  if (cn) {
    className.push(cn);
  }
  return <button {...props} className={className.join(' ')} style={style} type="button" />;
};

export default Button;
