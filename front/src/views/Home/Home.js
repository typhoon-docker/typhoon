import React, { useState, useEffect } from "react";

import { getProjects } from "/utils/axios";

import Project from "/components/Project/";

const Home = () => {
  const [projects, setProjects] = useState([]);
  const [selectedProjectID, selectProjectID] = useState(-1);

  useEffect(() => {
    getProjects().then(({ data }) => setProjects(data));
  }, []);

  const selectProject = id => event => {
    event.preventDefault();
    selectProjectID(id !== selectedProjectID ? id : -1);
  };

  return projects.map(project => (
    <Project
      key={project.id}
      project={project}
      onSelect={selectProject(project.id)}
      selected={project.id === selectedProjectID}
    />
  ));
};

export default Home;
