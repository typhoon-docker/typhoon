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
  const [MockAdapter, { default: randomArray }] = await Promise.all([
    import('axios-mock-adapter'),
    import('@libshin/random-array'),
  ]);
  const mock = new MockAdapter(client);

  const mockUser = await import('./mock/user.json');
  const mockProjects = (await import('./mock/projects.json')).map(project => ({ ...project, owner: mockUser }));
  const mockContainers = await import('./mock/containers.json');

  // getProjects
  mock.onGet('/projects').reply(200, mockProjects);

  // getAllProjects
  mock.onGet('/projects?all').reply(200, mockProjects);

  // postProject
  mock.onPost('/projects').reply(({ data }) => {
    const project = JSON.parse(data);
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
    const project = JSON.parse(data);
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

  // checkProject
  mock.onGet(/\/checkProject\?name=\w(\w|-)*/).reply(({ url, baseURL }) => {
    const name = url.substring(baseURL.length + 19);
    if (mockProjects.find(p => p.name === name)) {
      return [200, 'true'];
    }
    return [200, 'false'];
  });

  // statusProject
  mock.onGet(/\/docker\/status\/\w+/).reply(() => {
    return [200, randomArray(mockContainers, 3)];
  });

  // getLogs
  mock.onGet(/\/docker\/logs\/.*/).reply(() => {
    return [
      200,
      `Donec
vitae nisi vel
dolor pretium scelerisque.
Nulla mollis feugiat neque
eget convallis. Nulla laoreet urna enim, at
efficitur libero
placerat vitae. Suspendisse
potenti. Quisque sem lacus,
ultrices sed aliquam et, pellentesque quis
sapien. Integer nec
purus mollis, pretium
lectus id, tincidunt arcu.
Pellentesque aliquet risus a odio ultrices,
sed bibendum quam
venenatis. Fusce elementum
est viverra, cursus felis ne
uctus diam. Morbi vestibulum lorem
non nisl sagittis,
porttitor lobortis quam
interdum. Morbi ut sem a leo
tincidunt tincidunt eu vitae`,
    ];
  });
};

export const getProjects = () => client.get('/projects');
export const getAllProjects = () => client.get('/projects?all');
export const postProject = project => client.post('/projects', project);
export const getProject = projectID => client.get(`/projects/${projectID}`);
export const putProject = project => client.put(`/projects/${project.id}`, project);
export const deleteProject = projectID => client.delete(`/projects/${projectID}`);
export const checkProject = name => client.get(`/checkProject?name=${name}`);
export const activateProject = projectID => client.post(`/docker/apply/${projectID}`, {}, { timeout: 5 * 60 * 1000 });
export const startProject = projectID => client.post(`/docker/up/${projectID}`);
export const stopProject = projectID => client.post(`/docker/down/${projectID}`);
export const statusProject = projectID => client.get(`/docker/status/${projectID}`);
export const getLogs = (projectID, lines = 150) =>
  client.get(`/docker/logs/${projectID}?lines=${lines}`, { timeout: 1 * 60 * 1000 });
