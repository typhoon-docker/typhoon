import React from 'react';

import { getProject, putProject, deleteProject } from '/utils/typhoonAPI';

import useAxios from '/utils/useAxios';

import Box from '/components/Box';
import Logs from '/components/Logs';

const Project = ({ projectID }) => {
  const [project, setProject] = useAxios(getProject(projectID), {}, [projectID]);

  const onDelete = () => {
    if (window.confirm(`Tu es sûr de vouloir suppimer ${project.name}`)) {
      deleteProject(projectID);
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
