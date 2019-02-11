import React from 'react';

import { Input, Checkbox } from '/components/Input';

// external_domain_names: [],
// system_dependencies: [],
// dependency_files: [],
// databases: [],
// env: {},

const Variables = ({ project }) => {
  return (
    <>
      <h2 style={{ fontWeight: 600, padding: '1em', fontSize: '1.2em' }}>
        Paramètres du ton site{' '}
        <span style={{ fontWeight: 500, fontSize: '0.8em' }}>
          {"(aucun n'est obligatoire, tu peux les laisser vide)"}
        </span>
      </h2>

      <Input
        title="Dossier contenant le code (monorepo)"
        name="root_folder"
        placeholder="Exemple back"
        defaultValue={project ? project.root_folder : ''}
      />
      <Input title="Script d'installation" name="install_script" defaultValue={project ? project.install_script : ''} />
      <Input title="Script de build" name="build_script" defaultValue={project ? project.build_script : ''} />
      <Input title="Script de run" name="start_script" defaultValue={project ? project.start_script : ''} />
      <Input
        title="Dossier statique"
        placeholder="Exemple images"
        name="static_folder"
        defaultValue={project ? project.static_folder : ''}
      />
      <Input
        type="number"
        title="Port d'écoute"
        name="exposed_port"
        defaultValue={project ? project.exposed_port : ''}
        min={80}
        max={65535}
      />

      <Checkbox title="HTTPS" name="use_https" defaultChecked={project ? project.use_https : false} value="https" />
    </>
  );
};

export default Variables;
