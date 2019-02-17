import React from 'react';
import { navigate } from '@reach/router';

import { getProject, putProject, deleteProject } from '/utils/typhoonAPI';

import useAxios from '/utils/useAxios';

import Box from '/components/Box';
import Logs from '/components/Logs';

const Project = ({ projectID }) => {
  const [project, setProject] = useAxios(getProject(projectID), {}, [projectID]);

  const onDelete = () => {
    if (window.confirm(`Tu es sûr de vouloir suppimer ${project.name}`)) {
      deleteProject(projectID).then(() => navigate('/'));
    }
  };

  return (
    <Box>
      <Logs projectID={projectID} />
      <button type="button" onClick={onDelete}>
        Supprimer
      </button>
    </Box>
  );
};

export default Project;
