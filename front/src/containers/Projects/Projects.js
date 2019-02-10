import React, { useState, useEffect } from 'react';

import { getProjects, getAllProjects } from '/utils/axios';

import Box from '/components/Box/';
import Steps from '/components/Steps/';
import Project from '/components/Project/';
import EmptyProject from '/components/EmptyProject/';

import { title, button } from './Projects.css';

const Projects = ({ all = false }) => {
  const [projects, setProjects] = useState(null);
  const [selectedProjects, selectProject] = useState({});
  const [step, setStep] = useState(0);

  const fetchFn = all ? getAllProjects : getProjects;

  useEffect(() => {
    fetchFn().then(({ data }) => setProjects(data));
  }, []);

  const onSelect = id => event => {
    event.preventDefault();
    selectProject({ ...selectedProjects, [id]: !selectedProjects[id] });
  };

  return (
    <Steps step={step}>
      <Box>
        <h1 className={title}>
          {all ? 'Tous les projets' : 'Mes projets'}
          <button className={button} type="button" onClick={() => setStep(1)}>
            Nouveau projet
          </button>
        </h1>
        {projects &&
          (projects.length > 0 ? (
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
          ))}
      </Box>
    </Steps>
  );
};

export default Projects;