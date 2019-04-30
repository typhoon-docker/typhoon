import { createElement } from 'react';
import cx, { parseClassNames } from '/utils/className/';

const h = (as, { className, ...rest }) => {
  const { classNames, props } = parseClassNames(rest);
  return createElement(as, {
    className: cx(classNames, className),
    ...props,
  });
};

export default h;
