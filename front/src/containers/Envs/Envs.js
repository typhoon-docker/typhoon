import h from '/utils/h';

import { useState } from 'react';

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

const Envs = ({ defaultEnvs }) => {
  const [keys, setKeys] = useState(['']);
  const haveDefaultEnvs = defaultEnvs && Object.keys(defaultEnvs).length;

  if (keys.length === keys.filter(Boolean).length) {
    setKeys([...keys, '']);
  }

  const onChange = index => value => {
    const dump = [...keys];
    dump[index] = value;
    setKeys(dump);
  };
  return (
    <>
      {haveDefaultEnvs && (
        <>
          {Object.keys(defaultEnvs).map(defaultEnvName => (
            <input
              key={defaultEnvName}
              type="hidden"
              name={`env.${defaultEnvName}`}
              value={defaultEnvs[defaultEnvName]}
            />
          ))}
        </>
      )}
      {keys.map((key, index) => (
        // eslint-disable-next-line
        <Env key={index} value={key} onChange={onChange(index)} />
      ))}
    </>
  );
};

export default Envs;
