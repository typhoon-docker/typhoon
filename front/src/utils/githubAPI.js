import axios from 'axios';

import { newProjectCup } from './project';

let client;

const updateClient = token => {
  const config = {
    baseURL: 'https://api.github.com',
    headers: { Authorization: `token ${token}`, Accept: 'application/vnd.github.v3+json' },
  };
  client = axios.create(config);
};
newProjectCup.on(({ repository_token }) => updateClient(repository_token));
updateClient(newProjectCup());

export const importMocks = async () => {
  const MockAdapter = await import('axios-mock-adapter');
  const mock = new MockAdapter(client);

  const mockRepos = await import('./mock/repositories.json');
  const mockBranches = await import('./mock/branches.json');

  // getRepos
  mock.onGet('/user/repos').reply(200, mockRepos);

  // getBranches
  mock.onGet(/\/repos\/\w+\/\w+\/branches/).reply(200, mockBranches);
};

export const getRepos = () => client.get('/user/repos');
export const getBranches = project => client.get(`/repos/${project.full_name}/branches`);
