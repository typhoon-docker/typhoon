import axios from 'axios';

import { getRawToken, tokenCup } from './connect';

let client;

const updateClient = token => {
  const config = {
    baseURL: process.env.BACKEND_URL,
    headers: { Authorization: `Bearer ${token}` },
  };
  client = axios.create(config);
};
tokenCup.on(token => updateClient(token));
updateClient(getRawToken());

export const importMocks = async () => {
  const MockAdapter = await import('axios-mock-adapter');
  const mock = new MockAdapter(client);

  const mockUser = await import('./mock/user.json');
  const mockProjects = (await import('./mock/projects.json')).map(project => ({ ...project, owner: mockUser }));

  // getProjects
  mock.onGet('/projects').reply(200, mockProjects);

  // getAllProjects
  mock.onGet('/projects?all').reply(200, mockProjects);

  // postProject
  mock.onPost('/projects').reply(({ data }) => {
    const { project } = JSON.parse(data);
    if (!mockProjects.map(({ id }) => id).includes(project.id)) {
      mockProjects.push(project);
    }
    return [200, mockProjects];
  });

  // getProject
  mock.onGet(/\/projects\/\d+/).reply(({ url, baseURL }) => {
    const projectID = parseInt(url.substring(baseURL.length + 10), 10);
    const foundProject = mockProjects.find(project => project.id === projectID) || null;
    return [200, foundProject];
  });

  // putProject
  mock.onPut(/\/projects\/\d+/).reply(({ data, url, baseURL }) => {
    const projectID = parseInt(url.substring(baseURL.length + 10), 10);
    const { project } = JSON.parse(data);
    if (projectID === project.id && mockProjects.map(({ id }) => id).includes(projectID)) {
      const projectIndex = mockProjects.findIndex(p => p.id === projectID);
      mockProjects[projectIndex] = project;
    }
    return [200, mockProjects];
  });

  // deleteProject
  mock.onPut(/\/projects\/\d+/).reply(({ url, baseURL }) => {
    const projectID = parseInt(url.substring(baseURL.length + 10), 10);
    return [200, mockProjects.filter(project => project.id !== projectID)];
  });

  mock.onGet(/\/checkProject\?name=\w(\w|-)*/).reply(({ url, baseURL }) => {
    const name = url.substring(baseURL.length + 19);
    if (mockProjects.find(p => p.name === name)) {
      return [200, 'true'];
    }
    return [200, 'false'];
  });
};

export const getProjects = () => client.get('/projects');
export const getAllProjects = () => client.get('/projects?all');
export const postProject = project => client.post('/projects', project);
export const getProject = projectID => client.get(`/projects/${projectID}`);
export const putProject = project => client.put(`/projects/${project.id}`, project);
export const deleteProject = project => client.delete(`/projects/${project.id}`);
export const checkProject = name => client.get(`/checkProject?name=${name}`);
export const activateProject = projectID => client.get(`/docker/apply/${projectID}`);
export const startProject = projectID => client.get(`/docker/up/${projectID}`);
export const stopProject = projectID => client.get(`/docker/down/${projectID}`);
export const statusProject = projectID => client.get(`/docker/status/${projectID}`);
