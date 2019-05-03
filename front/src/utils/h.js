import { createElement, Fragment } from 'react';
import cx, { parseClassNames } from '/utils/className/';

const h = (as, p, ...c) => {
  if (!p) {
    return createElement(as, p, ...c);
  }
  const { className: cn, ...rest } = p;
  const { classNames, props } = parseClassNames(rest);
  const className = cx(classNames, cn);
  return createElement(
    as,
    className
      ? {
          className,
          ...props,
        }
      : props,
    ...c,
  );
};

h.Fragment = Fragment;

export default h;
