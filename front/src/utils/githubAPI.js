import axios from 'axios';
import { shouldMock } from '/utils/env';

import { newProjectCup } from './project';
import { arrayToJSON } from './formData';

let client;

export const importMocks = async () => {
  const [MockAdapter, { default: logMock }] = await Promise.all([import('axios-mock-adapter'), import('./mock/log')]);
  const mock = new MockAdapter(client);

  const [mockRepos, mockBranches, mockOrgs] = await Promise.all([
    import('./mock/repositories.json'),
    import('./mock/branches.json'),
    import('./mock/organizations.json'),
  ]);

  // getRepos
  mock.onGet('/user/repos').reply(logMock(200, mockRepos));

  // getBranches
  mock.onGet(/\/repos\/.*\/branches/).reply(logMock(200, mockBranches));

  // getOrgs
  mock.onGet('/user/orgs').reply(logMock(200, mockOrgs));

  // getOrgRepos
  mock.onGet('/orgs/orga/repos').reply(logMock(200, mockRepos));

  // Forward non mocked functions
  mock.onAny().passThrough();
};

const updateClient = async token => {
  const config = {
    baseURL: 'https://api.github.com',
    headers: { Authorization: `token ${token}`, Accept: 'application/vnd.github.v3+json' },
  };
  client = axios.create(config);
  if (shouldMock) {
    await importMocks();
  }
};

newProjectCup.on(({ repository_token }) => updateClient(repository_token));
updateClient(newProjectCup());

export const getRepos = async (page, admin = true) => {
  const res = await client.get(page ? `/user/repos?page=${page}` : `/user/repos`);
  if (admin) {
    res.data = res.data.filter(repo => repo.permissions.admin);
  }
  return res;
};
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
