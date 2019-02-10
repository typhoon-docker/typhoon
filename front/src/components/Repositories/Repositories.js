import React, { Fragment, useState, useEffect } from 'react';

import { getRepos } from '/utils/githubAPI';
import ArrowButton from '/components/ArrowButton/';

import { input, label } from './Repository.css';

const Repository = ({ repo }) => (
  <Fragment>
    <input type="radio" id={repo.id} className={input} name="repo" value={JSON.stringify(repo)} />
    <label htmlFor={repo.id} className={label}>
      {repo.name}
    </label>
  </Fragment>
);

const Repositories = () => {
  const [repos, setRepos] = useState([]);

  useEffect(() => {
    getRepos()
      .then(({ data }) => setRepos(data))
      .catch(console.warn);
  }, []);

  return (
    <div style={{ display: 'flex', flexDirection: 'column' }}>
      {repos.map(repo => (
        <Repository key={repo.id} repo={repo} />
      ))}
      <div style={{ marginTop: '0.5em', fontSize: '1.3em', alignSelf: 'flex-end' }}>
        <ArrowButton type="submit">Continuer</ArrowButton>
      </div>
    </div>
  );
};

export default Repositories;
