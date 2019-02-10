import React, { useState, useEffect } from 'react';

import { getRepos } from '/utils/githubAPI';

const Repositories = () => {
  const [repos, setRepos] = useState([]);

  useEffect(() => {
    getRepos()
      .then(({ data }) => setRepos(data))
      .catch(console.warn);
  }, []);

  return repos.map(repo => <pre>{JSON.stringify(repo, null, 2)}</pre>);
};

export default Repositories;
