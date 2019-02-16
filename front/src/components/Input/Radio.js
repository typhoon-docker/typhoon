import React from 'react';

import { label, radio } from './Input.css';
import { highlightable_ref, highlightable_offset } from '/styles/highlightable.css';

import cx from '/utils/className';

const Radio = ({ id, className: cn, children, ...props }) => (
  <>
    <input type="radio" id={id} className={cx(highlightable_ref, radio)} {...props} />
    <label htmlFor={id} className={cx(highlightable_offset, cn)}>
      {children}
    </label>
  </>
);

export default Radio;
