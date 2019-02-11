import { useState, useEffect } from 'react';
import { createCup } from 'manatea';

export const systemEnvVars = {
  TYPHOON_PERSISTENT_DIR: '/persistent',
};

export const newProjectCup = createCup({
  name: null,
  repository_url: null,
  repository_type: 'github',
  template_id: null,
  repository_token: null,
  external_domain_names: [],
  use_https: true,
  docker_image_version: null,
  root_folder: null,
  exposed_port: null,
  system_dependencies: [],
  dependency_files: [],
  install_script: null,
  build_script: null,
  start_script: null,
  static_folder: null,
  databases: [],
  env: {},
});

export const useProject = () => {
  const [project, setProject] = useState(newProjectCup());
  useEffect(() => {
    const listener = newProjectCup.on(p => setProject(p));
    setProject(newProjectCup());
    return listener;
  }, []);
  return [project, p => newProjectCup(p)];
};

export const isProjectFilled = () => {
  const project = newProjectCup();
  if (!project.name || !project.repository_url || !project.template_id || !project.repository_type) {
    return false;
  }
  return true;
};
