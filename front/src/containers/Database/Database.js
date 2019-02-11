import React, { useState } from 'react';

import { Checkbox, Input } from '/components/Input';
import list from '/utils/databases.json';

const Db = ({ title, update }) => (
  <>
    {title}
    <input type="hidden" name={`databases.${title}.tile`} value={title} placeholder="database" />
    <Input name={`databases.${title}.env_db`} required title="database" {...update('env_db')} placeholder="user" />
    <Input name={`databases.${title}.env_user`} required title="user" {...update('env_user')} />
    <Input
      name={`databases.${title}.env_password`}
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
          title={activeDb[index] ? <Db title={databases[index].type} update={update(index)} /> : databases[index].type}
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
