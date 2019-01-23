import React from 'react';

import { box } from './Box.css';

const Box = ({ className: cn, ...props }) => {
  const className = [box];
  if (cn) {
    className.push(cn);
  }
  return React.createElement('div', {
    ...props,
    className: className.join(' '),
  });
};

export default Box;
