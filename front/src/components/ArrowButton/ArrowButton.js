import React from 'react';

import { button, next, previous } from './ArrowButton.css';

import cx from '/utils/className';
import { useProperty } from '/utils/hooks';

const ArrowButton = ({ color = 'tertiary', direction = 'next', className: cn, as = 'button', ...props }) => {
  const ref = useProperty(
    {
      '--color': `var(--${color})`,
    },
    [color],
  );
  const className = [button, direction === 'previous' ? previous : next, cn];
  return React.createElement(as, {
    ref,
    types: 'button',
    ...props,
    className: cx(className),
  });
};

export default ArrowButton;
