import React from 'react';

import useProperty from '/utils/useProperty';
import { underline, bold } from './Underline.css';

const validStyle = [
  'none',
  'hidden',
  'dotted',
  'dashed',
  'solid',
  'double',
  'groove',
  'ridge',
  'inset',
  'outset',
  'initial',
  'inherit',
];

const Underline = ({
  as = 'span',
  color = 'text',
  className: cn,
  variant = 'solid',
  bold: isBold = false,
  style: s = {},
  ...props
}) => {
  const el = useProperty(
    () => ({
      '--color': `rgb(var(--${color}))`,
    }),
    [color],
  );

  let className = [underline];

  const style = {
    ...s,
    borderBottomStyle: variant && validStyle.includes(variant) ? variant : 'solid',
  };

  if (cn) {
    className.push(cn);
  }
  if (isBold) {
    className.push(bold);
  }

  className = className.join(' ');

  return React.createElement(as, {
    ...props,
    ref: el,
    className,
    style,
  });
};

export default Underline;
