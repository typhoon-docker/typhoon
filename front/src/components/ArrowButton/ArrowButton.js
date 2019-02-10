import React from 'react';

import { button } from './ArrowButton.css';

const ArrowButton = ({ className: cn, ...props }) => {
  const className = [button];
  if (cn) {
    className.push(cn);
  }
  return <button type="button" {...props} className={className.join(' ')} />;
};

export default ArrowButton;
