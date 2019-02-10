import React from 'react';
import murmur from '@emotion/hash';

import { wrapper } from './Input.css';

const Input = ({ title, error, errorMessage, ...props }) => {
  const id = murmur(`${Date.now()}_${title}`);
  return (
    <div className={wrapper}>
      <label htmlFor={id}>{title}</label>
      <input type="text" id={id} {...props} />
      {error && <span>{errorMessage}</span>}
    </div>
  );
};

export default Input;
