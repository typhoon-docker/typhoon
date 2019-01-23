import React, { useState, useEffect } from 'react';

import { getProjects } from '/utils/axios';

import Box from '/components/Box/';
import Project from '/components/Project/';
import EmptyProject from '/components/EmptyProject/';

import { title, button } from './Home.css';

const Home = () => {
  const [projects, setProjects] = useState(null);
  const [selectedProjects, selectProject] = useState({});

  useEffect(() => {
    getProjects().then(({ data }) => setProjects(data));
  }, []);

  const onSelect = id => event => {
    event.preventDefault();
    selectProject({ ...selectedProjects, [id]: !selectedProjects[id] });
  };

  return (
    <Box>
      <h1 className={title}>
        Mes projets
        <button className={button} type="button">
          Nouveau projet
        </button>
      </h1>
      {projects && projects.length > 0 ? (
        projects.map(project => (
          <Project
            key={project.id}
            project={project}
            onSelect={onSelect(project.id)}
            selected={selectedProjects[project.id]}
          />
        ))
      ) : (
        <EmptyProject />
      )}
    </Box>
  );
};

export default Home;
