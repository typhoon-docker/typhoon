import React from 'react';

import injectClassName from '/utils/injectClassName';

import { wrapper, hidden_left, visible_middle, hidden_right } from './Steps.css';

const Steps = ({ step = 0, children: c, ...props }) => {
  const children = React.Children.toArray(c);
  const nbChildren = children.length;
  if (nbChildren === 0) {
    return null;
  }

  let index;
  if (!(step in children)) {
    if (step < 0) {
      index = 0;
    } else if (step >= nbChildren) {
      index = nbChildren - 1;
    } else if (parseInt(step, 0) in children) {
      index = parseInt(step, 0);
    } else {
      return null;
    }
  } else {
    index = step;
  }
  return (
    <div {...props} className={wrapper}>
      {index - 1 in children && injectClassName(children[index - 1], hidden_left, { key: index - 1 })}
      {injectClassName(children[index], visible_middle, { key: index })}
      {index + 1 in children && injectClassName(children[index + 1], hidden_right, { key: index + 1 })}
    </div>
  );
};

export default Steps;
