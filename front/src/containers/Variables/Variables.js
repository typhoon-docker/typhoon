import React from 'react';

import { Input, Checkbox } from '/components/Input';
import globalIgnoreField from '/utils/ignore_fields.json';

const Variables = ({ project }) => {
  if (!project) {
    return null;
  }
  const {
    template_id,
    external_domain_names,
    root_folder,
    system_dependencies,
    dependency_files,
    install_script,
    build_script,
    start_script,
    static_folder,
    exposed_port,
    use_https,
  } = project;

  const ignoreField = template_id ? globalIgnoreField[template_id] : {};

  const questions = [
    {
      title: 'Noms de domaines externes',
      name: 'external_domain_names',
      placeholder: 'Exemple : mon-site.fr,www.mon-site.fr (séparés par une virgule, sans https ni / à la fin)',
      agreement: 'plural',
      defaultValue: external_domain_names ? external_domain_names.join(',') : '',
    },
    {
      title: 'Dossier contenant le code (monorepo)',
      name: 'root_folder',
      placeholder: 'Exemple back',
      defaultValue: root_folder || '',
    },
    {
      title: 'Dépendances systèmes',
      name: 'system_dependencies',
      placeholder: 'Exemple : ffmpeg,imagemagick (séparés par une virgule)',
      agreement: 'plural',
      defaultValue: system_dependencies ? system_dependencies.join(',') : '',
    },
    {
      title: 'Fichiers de dépendances',
      name: 'dependency_files',
      placeholder: 'Exemple : package.json,yarn.lock (séparés par une virgule)',
      agreement: 'plural',
      defaultValue: dependency_files ? dependency_files.join(',') : '',
    },
    {
      title: "Script d'installation",
      name: 'install_script',
      defaultValue: install_script || '',
    },
    {
      title: 'Script de build',
      name: 'build_script',
      defaultValue: build_script || '',
    },
    {
      title: 'Script de run',
      name: 'start_script',
      defaultValue: start_script || '',
    },
    {
      title: 'Dossier statique',
      placeholder: 'Exemple images',
      name: 'static_folder',
      defaultValue: static_folder || '',
    },
    {
      type: 'number',
      title: "Port d'écoute",
      name: 'exposed_port',
      defaultValue: exposed_port || '',
      min: 80,
      max: 65535,
    },
  ];
  return (
    <>
      <h2 style={{ fontWeight: 600, padding: '1em', fontSize: '1.2em' }}>
        Paramètres du ton site{' '}
        <span style={{ fontWeight: 500, fontSize: '0.8em' }}>
          {"(aucun n'est obligatoire, tu peux les laisser vide)"}
        </span>
      </h2>

      {questions.map(question =>
        ignoreField[question.name] ? null : <Input key={question.name} askIfEmpty {...question} />,
      )}

      <Checkbox title="HTTPS" name="use_https" defaultChecked={project ? use_https : false} value="https" />
    </>
  );
};

export default Variables;
