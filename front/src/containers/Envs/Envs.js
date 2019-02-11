import React, { useState } from 'react';

import { input } from './Envs.css';

const Env = ({ value, onChange }) => {
  return (
    <div style={{ display: 'flex' }}>
      <input
        type="text"
        className={input}
        value={value}
        onChange={event => onChange(event.target.value.toUpperCase().replace(/\W+/g, '_'))}
        placeholder="Variable d'environnement"
        style={{ flex: 1 }}
      />
      <input
        type="text"
        className={input}
        {...(value ? { name: `env.${value}` } : {})}
        placeholder="Valeur"
        style={{ flex: 2 }}
      />
    </div>
  );
};

const Envs = () => {
  const [keys, setKeys] = useState(['']);

  if (keys.length === keys.filter(Boolean).length) {
    setKeys([...keys, '']);
  }

  const onChange = index => value => {
    const dump = [...keys];
    dump[index] = value;
    setKeys(dump);
  };

  // eslint-disable-next-line
  return keys.map((key, index) => <Env key={index} value={key} onChange={onChange(index)} />);
};

export default Envs;
