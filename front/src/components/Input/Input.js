import React, { useState } from 'react';
import murmur from '@emotion/hash';

import { wrapper, label, error as errorCN, input } from './Input.css';
import { auto_highlightable as highlightable } from '/styles/highlightable.css';

import Button from '/components/Button';

import cx from '/utils/className';

const Input = ({ title, error, defaultValue, askIfEmpty, errorMessage, ...props }) => {
  const [showInput, setShowInput] = useState(false);

  const id = murmur(`${Date.now()}_${title}`);

  return defaultValue || !askIfEmpty || showInput ? (
    <div className={wrapper}>
      <label htmlFor={id} className={label}>
        {title}
      </label>
      <input type="text" {...props} id={id} defaultValue={defaultValue} className={cx(highlightable, input)} />
      {error && <span className={errorCN}>{errorMessage}</span>}
    </div>
  ) : (
    <div style={{ color: 'var(--text)' }}>
      <Button color="rgba(var(--secondary), 0.3)" type="button" onClick={() => setShowInput(true)}>
        <span role="img" aria-label="plus">
          ➕
        </span>
        Ajouter un <span style={{ textTransform: 'lowercase' }}>{title}</span>
      </Button>
    </div>
  );
};

export default Input;
