import React, { Fragment, useState, useEffect } from 'react';

import { getRepos } from '/utils/githubAPI';

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
      {repo.full_name}
    </label>
  </Fragment>
);

const Repositories = ({ onSelect }) => {
  const [repositories, setRepos] = useState([]);

  const fetchRepos = page =>
    getRepos(page)
      .then(({ data, headers }) => {
        setRepos(r => [...r, ...data]);
        if (headers.link && headers.link.includes(';')) {
          fetchRepos(
            parseInt(
              headers.link
                .split(';')[0]
                .replace(/[<>]/g, '')
                .split('?page=')[1],
              10,
            ),
          );
        }
      })
      .catch(console.warn);

  useEffect(() => {
    fetchRepos(0);
  }, []);

  return repositories.map(repo => <Repository key={repo.id} repo={repo} onSelect={onSelect} />);
};

export default Repositories;
