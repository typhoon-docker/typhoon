import React, { useEffect, useRef } from 'react';
import Color from 'color';

import { button } from './Button.css';

const regex = /^[a-z]([a-z]|-)*$/;

const Button = ({ color: c, className: cn, style: s, ...props }) => {
  const ref = useRef(null);

  const style = s ? { ...s } : {};
  const color = c || 'text';
  style.backgroundColor = color.match(regex) ? `rgb(var(--${color}))` : color;

  const className = [button];
  if (cn) {
    className.push(cn);
  }

  useEffect(() => {
    const textColor = window.getComputedStyle(ref.current).backgroundColor;

    if (Color(textColor).isLight()) {
      ref.current.style.color = 'rgb(var(--text))';
    } else {
      ref.current.style.color = 'white';
    }
  });

  return <button type="button" {...props} className={className.join(' ')} style={style} ref={ref} />;
};

export default Button;
