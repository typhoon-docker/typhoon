import { button, next, previous } from './ArrowButton.css';

import h from '/utils/h';
import { useProperty } from '/utils/hooks';

const ArrowButton = ({ color = 'tertiary', direction = 'next', className: cn, as = 'button', ...props }) => {
  const ref = useProperty(
    {
      '--color': `var(--${color})`,
    },
    [color],
  );
  const className = [button, direction === 'previous' ? previous : next, cn];
  return h(as, {
    ref,
    types: 'button',
    ...props,
    className,
  });
};

export default ArrowButton;
