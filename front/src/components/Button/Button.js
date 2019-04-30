import h from '/utils/h';

import { button } from './Button.css';

import { auto_highlightable as highlightable } from '/styles/highlightable.css';
import { useProperty } from '/utils/hooks';

const Button = ({ color, className: cn, style: s, ...props }) => {
  const ref = useProperty(
    {
      '--color': `var(--${color || 'text'})`,
    },
    [color],
  );

  const style = s ? { ...s } : {};

  return h('button', {
    type: 'button',
    ...props,
    className: [button, highlightable, cn],
    style,
    ref,
  });
};

export default Button;
