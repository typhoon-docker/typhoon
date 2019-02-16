import React from 'react';

import { button, next, previous } from './ArrowButton.css';

import cx from '/utils/className';

const ArrowButton = ({ direction = 'next', className: cn, ...props }) => {
  const className = [button, direction === 'previous' ? previous : next, cn];
  return <button type="button" {...props} className={cx(className)} />;
};

export default ArrowButton;
