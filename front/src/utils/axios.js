import axios from 'axios';

import { getRawToken } from './connect';

const config = {
  baseURL: process.env.BACKEND_URL,
  headers: { Authorization: `Bearer ${getRawToken()}` },
};

const client = axios.create(config);

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
  mock.onGet(/\/projects\/\d+/).reply(({ url }) => {
    const projectID = parseInt(url.substring(10), 10);
    const foundProject = mockProjects.find(project => project.id === projectID) || null;
    return [200, foundProject];
  });

  // putProject
  mock.onPut(/\/projects\/\d+/).reply(({ data, url }) => {
    const projectID = parseInt(url.substring(10), 10);
    const { project } = JSON.parse(data);
    if (projectID === project.id && mockProjects.map(({ id }) => id).includes(projectID)) {
      const projectIndex = mockProjects.findIndex(p => p.id === projectID);
      mockProjects[projectIndex] = project;
    }
    return [200, mockProjects];
  });

  // deleteProject
  mock.onPut(/\/projects\/\d+/).reply(({ url }) => {
    const projectID = parseInt(url.substring(10), 10);
    return [200, mockProjects.filter(project => project.id !== projectID)];
  });
};

export const getProjects = () => client.get('/projects');
export const getAllProjects = () => client.get('/projects?all');
export const postProject = project => client.post('/projects', { project });
export const getProject = projectID => client.get(`/projects/${projectID}`);
export const putProject = project => client.put(`/projects/${project.id}`, { project });
export const deleteProject = project => client.delete(`/projects/${project.id}`);
