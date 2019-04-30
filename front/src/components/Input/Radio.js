import React from 'react';

import { content, radio, auto_check } from './Input.css';
import { highlightable_ref, highlightable_offset } from '/styles/highlightable.css';

import cx from '/utils/className/';

const Radio = ({ autoCheck, id, className: cn, children, ...props }) => (
  <>
    <input type="radio" id={id} className={cx(highlightable_ref, radio, autoCheck && auto_check)} {...props} />
    <label htmlFor={id} className={cx(highlightable_offset, content, cn)}>
      {children}
    </label>
  </>
);

export default Radio;
