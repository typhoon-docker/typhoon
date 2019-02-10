import React from 'react';

import { box } from './Box.css';

const Box = ({ className: cn, as = 'div', ...props }) => {
  const className = [box];
  if (cn) {
    className.push(cn);
  }
  return React.createElement(as, {
    ...props,
    className: className.join(' '),
  });
};

export default Box;
