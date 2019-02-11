import React, { useState } from 'react';

import { Checkbox, Input } from '/components/Input';
import list from '/utils/databases.json';

const Db = ({ db, update }) => (
  <>
    {db.type}
    {Object.entries(db)
      .filter(entries => entries[1] !== '')
      .map(([key, value]) => (
        <input type="hidden" key={key} name={`databases.${db.type}.${key}`} value={value} />
      ))}
    <Input
      name={`databases.${db.type}.env_db`}
      required
      title="database"
      {...update('env_db')}
      placeholder="database"
    />
    <Input name={`databases.${db.type}.env_user`} required title="user" {...update('env_user')} placeholder="user" />
    <Input
      name={`databases.${db.type}.env_password`}
      required
      title="password"
      {...update('env_password')}
      type="password"
      placeholder="password"
    />
  </>
);

const Database = () => {
  const [activeDb, setActiveDbs] = useState(new Array(list.length).fill(false));
  const [databases, setDatabases] = useState(list);

  const update = index => field => ({
    value: databases[index][field],
    onChange: event => {
      const dump = [...databases];
      dump[index] = { ...dump[index], [field]: event.target.value };
      setDatabases(dump);
    },
  });

  return (
    <>
      {databases.map((db, index) => (
        <Checkbox
          key={db.type}
          title={activeDb[index] ? <Db db={db} update={update(index)} /> : db.type}
          checked={activeDb[index]}
          onChange={event => {
            const dump = [...activeDb];
            dump[index] = event.target.checked;
            setActiveDbs(dump);
          }}
        />
      ))}
    </>
  );
};

export default Database;
