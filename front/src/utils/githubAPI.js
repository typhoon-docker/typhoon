import axios from 'axios';

import { newProjectCup } from './project';
import { arrayToJSON } from './formData';

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
  const mockOrgs = await import('./mock/organizations.json');

  // getRepos
  mock.onGet('/user/repos').reply(200, mockRepos);

  // getBranches
  mock.onGet(/\/repos\/\w+\/\w+\/branches/).reply(200, mockBranches);

  // getOrgs
  mock.onGet('/user/orgs').reply(200, mockOrgs);

  // getOrgRepos
  mock.onGet('/orgs/orga/repos').reply(200, mockRepos);
};

export const getRepos = (page = 0) => client.get(`/user/repos?page=${page}`);
export const getBranches = project => client.get(`/repos/${project.full_name}/branches`);
export const getOrgs = () => client.get('/user/orgs');
export const getOrgRepos = () =>
  client
    .get('/user/orgs')
    .then(({ data }) => Promise.all(data.map(async org => [org.login, (await client.get(org.repos_url)).data])))
    .then(arrayToJSON);
export const getFullRepos = async () => {
  const [self, orgs] = await Promise.all([getRepos(), getOrgRepos()]);
  return { data: { 'Vos répos': self.data, ...orgs } };
};
