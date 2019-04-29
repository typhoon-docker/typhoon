import React, { useState } from 'react';

import { getProjects, getAllProjects } from '/utils/typhoonAPI';
import { useAxios } from '/utils/hooks';

import Box from '/components/Box/';
import ArrowButton from '/components/ArrowButton/';
import Button from '/components/Button/';
import Steps from '/components/Steps/';
import Project from '/components/Project/';
import EmptyProject from '/components/EmptyProject/';

import { title, git_wrapper } from './Projects.css';

const Projects = ({ all = false }) => {
  const fetchFn = all ? getAllProjects : getProjects;

  const [projects] = useAxios(fetchFn, null, []);
  const [selectedProjects, selectProject] = useState({});
  const [step, setStep] = useState(0);

  const onSelect = id => event => {
    event.preventDefault();
    selectProject({ ...selectedProjects, [id]: !selectedProjects[id] });
  };

  return (
    <Steps step={step}>
      <Box>
        <h1 className={title}>
          {all ? 'Tous les projets' : 'Mes projets'}
          <ArrowButton type="button" onClick={() => setStep(1)}>
            Nouveau projet
          </ArrowButton>
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
      <Box>
        <h2 className={title}>Quel est votre gestionnaire de Code ?</h2>
        <div className={git_wrapper}>
          <Button
            color="github"
            onClick={() => {
              window.location.href = `${process.env.BACKEND_URL}/login/github`;
            }}
          >
            GitHub
          </Button>
        </div>
      </Box>
    </Steps>
  );
};

export default Projects;
