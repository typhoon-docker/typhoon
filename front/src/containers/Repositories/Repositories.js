import React, { Fragment, useState, useEffect } from 'react';

import { getFullRepos } from '/utils/githubAPI';

import { input, label } from './Repository.css';

const Repository = ({ id, repo, onSelect }) => (
  <Fragment>
    <input
      type="radio"
      id={id}
      className={input}
      name="repository_url"
      value={repo.url}
      onChange={() => onSelect(repo)}
    />
    <label htmlFor={id} className={label}>
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

  return Object.entries(repositories).map(([org, repos]) => (
    <Fragment key={org}>
      <h2>{org}</h2>
      {repos.map(repo => (
        <Repository key={`${org}_${repo.id}`} id={`${org}_${repo.id}`} repo={repo} onSelect={onSelect} />
      ))}
    </Fragment>
  ));
};

export default Repositories;
