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

  useEffect(() => {
    const fetchRepos = page =>
      getRepos(page)
        .then(({ data, headers }) => {
          setRepos(r => [...r, ...data]);
          if (headers.link && headers.link.includes(';')) {
            const nextPage = Number(headers.link.split('>; rel="next"')[0].replace(/.*\?page=/, ''));
            if (!Number.isNaN(nextPage)) {
              fetchRepos(nextPage);
            }
          }
        })
        .catch(console.warn);
    fetchRepos();
  }, []);

  return repositories.map(repo => <Repository key={repo.id} repo={repo} onSelect={onSelect} />);
};

export default Repositories;
