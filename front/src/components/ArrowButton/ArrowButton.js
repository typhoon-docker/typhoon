import React from 'react';

import { button, next, previous } from './ArrowButton.css';

const ArrowButton = ({ direction = 'next', className: cn, ...props }) => {
  const className = [button, direction === 'previous' ? previous : next];
  if (cn) {
    className.push(cn);
  }
  return <button type="button" {...props} className={className.join(' ')} />;
};

export default ArrowButton;
