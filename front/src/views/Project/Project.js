import React from 'react';
import { navigate } from '@reach/router';

import { getProject, activateProject, putProject, deleteProject } from '/utils/typhoonAPI';

import useAxios from '/utils/useAxios';

import Box from '/components/Box';
import Logs from '/components/Logs';
import { BuildLogs, DockerFiles } from '/components/LinesBlocks';

const Project = ({ projectID }) => {
  const [project, setProject] = useAxios(getProject(projectID), {}, [projectID]);
  // ToDo: Use setProject and putProject to modify the project

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
      <BuildLogs projectID={projectID} />
      <DockerFiles projectID={projectID} />
    </Box>
  );
};

export default Project;
