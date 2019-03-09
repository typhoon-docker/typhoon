import React from 'react';
import { navigate } from '@reach/router';

import { getProject, putProject, deleteProject } from '/utils/typhoonAPI';

import useAxios from '/utils/useAxios';

import Box from '/components/Box';
import Logs from '/components/Logs';

const Project = ({ projectID }) => {
  const [project, setProject] = useAxios(getProject(projectID), {}, [projectID]);

  const onRedeploy = () => {
    if (window.confirm(`Le projet "${project.name}" va être redéployé`)) {
      activateProject(projectID);
    }
  };

  const onDelete = () => {
    if (window.confirm(`Tu es sûr de vouloir suppimer ${project.name}`)) {
      deleteProject(projectID).then(() => navigate('/'));
    }
  };

  return (
    <Box>
      <pre>{JSON.stringify(project, null, 2)}</pre>
      <button type="button" onClick={onRedeploy}>
        Redéployer
      </button>
      <button type="button" onClick={onDelete}>
        Supprimer
      </button>
      <Logs projectID={projectID} />
    </Box>
  );
};

export default Project;
