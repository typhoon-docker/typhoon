import React from 'react';
import murmur from '@emotion/hash';

import { wrapper, label, error as errorCN } from './Input.css';

const Select = ({ title, data, name, error, errorMessage, ...props }) => {
  const id = murmur(`${Date.now()}_${title}`);
  return (
    <div className={wrapper}>
      <label htmlFor={id} className={label}>
        {title}
      </label>
      <select name={name} id={id} {...props}>
        {data.map(({ key, value, display }) => (
          <option key={key || value} value={value}>
            {display || value}
          </option>
        ))}
      </select>
      {error && <span className={errorCN}>{errorMessage}</span>}
    </div>
  );
};

export default Select;
