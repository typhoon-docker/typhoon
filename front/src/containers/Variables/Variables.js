import React from 'react';

import { Input, Checkbox } from '/components/Input';

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
        title="Noms de domaines externes"
        name="external_domain_names"
        placeholder="Exemple : https://mon-site.fr/,https://www.mon-site.fr/ (séparés par une virgule)"
        askIfEmpty
        agreement="plural"
        defaultValue={project ? project.external_domain_names.join(',') : ''}
      />

      <Input
        title="Dossier contenant le code (monorepo)"
        name="root_folder"
        placeholder="Exemple back"
        askIfEmpty
        defaultValue={project ? project.root_folder : ''}
      />

      <Input
        title="Dépendances systèmes"
        name="system_dependencies"
        placeholder="Exemple : ffmpeg,imagemagick (séparés par une virgule)"
        askIfEmpty
        agreement="plural"
        defaultValue={project ? project.system_dependencies.join(',') : ''}
      />
      <Input
        title="Fichiers de dépendances"
        name="dependency_files"
        placeholder="Exemple : package.json,yarn.lock (séparés par une virgule)"
        askIfEmpty
        agreement="plural"
        defaultValue={project ? project.dependency_files.join(',') : ''}
      />
      <Input
        title="Script d'installation"
        name="install_script"
        askIfEmpty
        defaultValue={project ? project.install_script : ''}
      />

      <Input
        title="Script de build"
        name="build_script"
        askIfEmpty
        defaultValue={project ? project.build_script : ''}
      />
      <Input title="Script de run" name="start_script" askIfEmpty defaultValue={project ? project.start_script : ''} />

      <Input
        title="Dossier statique"
        placeholder="Exemple images"
        name="static_folder"
        askIfEmpty
        defaultValue={project ? project.static_folder : ''}
      />

      <Input
        type="number"
        title="Port d'écoute"
        name="exposed_port"
        askIfEmpty
        defaultValue={project ? project.exposed_port : ''}
        min={80}
        max={65535}
      />
      <Checkbox title="HTTPS" name="use_https" defaultChecked={project ? project.use_https : false} value="https" />
    </>
  );
};

export default Variables;
