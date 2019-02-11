import React from 'react';
import murmur from '@emotion/hash';

import { wrapper, label, input, error as errorCN } from './Input.css';

const Input = ({ title, error, errorMessage, ...props }) => {
  const id = murmur(`${Date.now()}_${title}`);
  return (
    <div className={wrapper}>
      <label htmlFor={id} className={label}>
        {title}
      </label>
      <input type="text" {...props} id={id} className={input} />
      {error && <span className={errorCN}>{errorMessage}</span>}
    </div>
  );
};

export default Input;
