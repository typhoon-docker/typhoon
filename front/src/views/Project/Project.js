import h from '/utils/h';
import { navigate, Link } from '@reach/router';

import { getProjects, getProject, activateProject, putProject, deleteProject } from '/utils/typhoonAPI';

import { useAxios } from '/utils/hooks';

import useFullBody from '/utils/useFullBody';

// import Box from '/components/Box';
import Logs from '/components/Logs';
import { BuildLogs, DockerFiles } from '/components/LinesBlocks';

import { container, projects_list, project_item, project_details, env_section, build_section } from './Project.css';

const Project = ({ projectID }) => {
  const [project, setProject] = useAxios(() => getProject(projectID), {}, [projectID]);
  const [projects] = useAxios(() => getProjects(), [], []);
  // ToDo: Use setProject and putProject to modify the project

  useFullBody();

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
    <div className={container}>
      <div className={projects_list}>
        {projects.map(p => (
          <Link key={p.id} to={`/project/${p.id}`} className={project_item} $class-active={project.id === p.id}>
            {p.name}
          </Link>
        ))}
      </div>
      <div className={project_details}>
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
      </div>
    </div>
  );
};

export default Project;
