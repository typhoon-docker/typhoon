import h from '/utils/h';
import { forwardRef } from 'react';

import { box } from './Box.css';

const Box = ({ className: cn, as = 'div', ...props }, ref) => {
  const className = [box];
  if (cn) {
    className.push(cn);
  }
  return h(as, {
    ...props,
    ref,
    className: className.join(' '),
  });
};

export default forwardRef(Box);
