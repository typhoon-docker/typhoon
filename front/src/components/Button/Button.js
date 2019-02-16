import React, { useRef, useEffect } from 'react';

import { button } from './Button.css';

import { auto_highlightable as highlightable } from '/styles/highlightable.css';
import cx from '/utils/className';

const Button = ({ color, className: cn, style: s, ...props }) => {
  const ref = useRef(null);

  const style = s ? { ...s } : {};

  useEffect(() => {
    ref.current.style.setProperty('--color', `var(--${color || 'text'})`);
  });

  return <button type="button" {...props} className={cx(button, highlightable, cn)} style={style} ref={ref} />;
};

export default Button;
