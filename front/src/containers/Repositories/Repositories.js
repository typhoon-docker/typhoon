import h from '/utils/h';

import { useState, useEffect } from 'react';

import { getRepos } from '/utils/githubAPI';

import { Radio } from '/components/Input';

const Repositories = ({ onSelect }) => {
  const [repositories, setRepos] = useState([]);

  useEffect(() => {
    const fetchRepos = page =>
      getRepos(page)
        .then(({ data, headers }) => {
          setRepos(r => [...r, ...data]);
          if (headers && headers.link && headers.link.includes(';')) {
            const nextPage = Number(headers.link.split('>; rel="next"')[0].replace(/.*\?page=/, ''));
            if (!Number.isNaN(nextPage)) {
              fetchRepos(nextPage);
            }
          }
        })
        .catch(console.warn);
    fetchRepos();
  }, []);

  return repositories.map(repo => (
    <Radio
      key={repo.id}
      id={repo.id}
      name="repository_url"
      value={repo.clone_url}
      onChange={() => onSelect(repo)}
      autoCheck
    >
      {repo.full_name}
    </Radio>
  ));
};

export default Repositories;
