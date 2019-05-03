import h from '/utils/h';

import murmur from '@emotion/hash';

import { label, input, error as errorCN } from './Input.css';

const Checkbox = ({ title, error, errorMessage, ...props }) => {
  const id = murmur(`${Date.now()}_${title}`);
  return (
    <div>
      <input type="checkbox" {...props} id={id} className={input} />
      <label htmlFor={id} className={label}>
        {title}
      </label>
      {error && <span className={errorCN}>{errorMessage}</span>}
    </div>
  );
};

export default Checkbox;
