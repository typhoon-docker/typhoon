import React, { useState } from 'react';
import murmur from '@emotion/hash';

import { wrapper, label, error as errorCN, input } from './Input.css';
import { auto_highlightable as highlightable } from '/styles/highlightable.css';

import Button from '/components/Button';

import cx from '/utils/className/';

const agreements = agreement => {
  if (agreement === 'plural') {
    return 'des';
  }
  if (agreement === 'feminine') {
    return 'une';
  }
  return 'un';
};

const Input = ({ title, error, agreement, defaultValue, askIfEmpty, errorMessage, ...props }) => {
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
      <Button color="primary" className="highlightable" type="button" onClick={() => setShowInput(true)}>
        <span role="img" aria-label="plus">
          ➕
        </span>
        Ajouter {agreements(agreement)} <span style={{ textTransform: 'lowercase' }}>{title}</span>
      </Button>
    </div>
  );
};

export default Input;
