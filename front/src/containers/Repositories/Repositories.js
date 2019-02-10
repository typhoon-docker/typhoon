import React, { Fragment, useState, useEffect } from 'react';

import { getFullRepos } from '/utils/githubAPI';
import ArrowButton from '/components/ArrowButton/';

import { input, label } from './Repository.css';

const Repository = ({ repo, onSelect }) => (
  <Fragment>
    <input
      type="radio"
      id={repo.id}
      className={input}
      name="repository_url"
      value={repo.url}
      onChange={() => onSelect(repo)}
    />
    <label htmlFor={repo.id} className={label}>
      {repo.name}
    </label>
  </Fragment>
);

const Repositories = ({ onSelect }) => {
  const [repositories, setRepos] = useState({});

  useEffect(() => {
    getFullRepos()
      .then(({ data }) => setRepos(data))
      .catch(console.warn);
  }, []);

  return (
    <div style={{ display: 'flex', flexDirection: 'column' }}>
      {Object.entries(repositories).map(([org, repos]) => (
        <Fragment key={org}>
          <h2>{org}</h2>
          {repos.map(repo => (
            <Repository key={repo.id} repo={repo} onSelect={onSelect} />
          ))}
        </Fragment>
      ))}
      <div style={{ marginTop: '0.5em', fontSize: '1.3em', alignSelf: 'flex-end' }}>
        <ArrowButton type="submit">Continuer</ArrowButton>
      </div>
    </div>
  );
};

export default Repositories;
