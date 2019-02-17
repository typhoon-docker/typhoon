import React from 'react';

import { button } from './Button.css';

import { auto_highlightable as highlightable } from '/styles/highlightable.css';
import cx from '/utils/className';
import useProperty from '/utils/useProperty';

const Button = ({ color, className: cn, style: s, ...props }) => {
  const ref = useProperty(
    () => ({
      '--color': `var(--${color || 'text'})`,
    }),
    [color],
  );

  const style = s ? { ...s } : {};

  return <button type="button" {...props} className={cx(button, highlightable, cn)} style={style} ref={ref} />;
};

export default Button;
