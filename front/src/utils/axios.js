import axios from "axios";

import { getRawToken } from "./connect";
import { shouldMock } from "./env";

const config = {
  baseURL: process.env.BACKEND_URL,
  headers: { Authorization: `Bearer ${getRawToken()}` }
};

const client = axios.create(config);

if (shouldMock) {
  import("axios-mock-adapter").then(async MockAdapter => {
    const mock = new MockAdapter(client);

    const mockUser = await import("./mock/user.json");
    const mockProjects = (await import("./mock/projects.json")).map(project => ({ ...project, owner: mockUser }));

    mock.onGet("/projects").reply(200, mockProjects);
    mock.onGet("/projects?all").reply(200, mockProjects);
  });
}

export const getProjects = () => client.get("/projects");

export const getAllProjects = () => client.get("/projects?all");
