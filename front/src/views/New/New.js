import React, { useState, useEffect } from 'react';
import { useInfuser } from 'react-manatea';
import { Redirect } from '@reach/router';

import Repositories from '/containers/Repositories/';
import Database from '/containers/Database/';
import Envs from '/containers/Envs/';

import Steps from '/components/Steps/';
import Box from '/components/Box/';
import { Input, Select } from '/components/Input/';
import ArrowButton from '/components/ArrowButton/';
import TemplatePicker from '/components/TemplatePicker/';
import Variables from '/components/Variables/';

import { newProjectCup } from '/utils/project';
import { formDataToArray, arrayToJSON } from '/utils/formData';
import { checkProject, postProject, activateProject } from '/utils/typhoonAPI';
import { getBranches } from '/utils/githubAPI';
import globalIgnoreField from '/utils/ignore_fields.json';

import { block, direction } from './New.css';

const split = string =>
  (string || '')
    .split(',')
    .map(s => s.trim())
    .filter(Boolean);

const New = () => {
  const [step, setStep] = useState(0);
  const [error, setError] = useState(null);
  const [repo, setRepo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [template, setTemplate] = useState({});
  const [branches, setBranches] = useState([]);
  const [project, setProject] = useInfuser(newProjectCup);

  useEffect(() => {
    if (repo) {
      getBranches(repo)
        .then(({ data }) => {
          setBranches(data);
        })
        .catch(console.warn);
    }
  }, [repo]);

  if (!project.repository_token && process.env.NODE_ENV !== 'development') {
    return <Redirect to="/" />;
  }

  const nextStep = () => setStep(step + 1);
  const previousStep = () => setStep(step - 1);

  const onSubmit = (verifies = {}, transform, cb = () => {}) => event => {
    event.preventDefault();
    const newDataArray = formDataToArray(new FormData(event.target));
    if (Object.keys(verifies).length > newDataArray.length) {
      return;
    }
    if (Object.keys(verifies).length < newDataArray.length) {
      console.warn("Many fields aren't verified");
    }
    Promise.all(newDataArray.map(([key, value]) => (verifies[key] ? verifies[key](value) : true)))
      .then(array => {
        const falseIndex = array.findIndex(e => e === false);
        if (falseIndex >= 0) {
          setError(newDataArray[falseIndex][0]);
          throw new Error(`Error in ${newDataArray[falseIndex]}`);
        }
        return array;
      })
      .then(array => array.reduce((acc, isValid, index) => (isValid ? [...acc, newDataArray[index]] : acc), []))
      .then(array => {
        setProject(p => {
          let newData = arrayToJSON(array);
          newData = { ...newData, ...(transform || (x => x))(newData, p) };

          return { ...p, ...newData };
        });
        nextStep();
      })
      .then(cb)
      .catch(console.warn);
  };

  const ignoreField = globalIgnoreField[project.template_id] || {};
  const questions = [
    {
      title: 'Dossier contenant le code (monorepo)',
      name: 'root_folder',
      placeholder: 'Exemple back',
      defaultValue: project.root_folder || '',
    },
    {
      title: 'Dépendances systèmes',
      name: 'system_dependencies',
      placeholder: 'Exemple : ffmpeg,imagemagick (séparés par une virgule)',
      agreement: 'plural',
      defaultValue: project.system_dependencies ? project.system_dependencies.join(',') : '',
    },
    {
      title: 'Fichiers de dépendances',
      name: 'dependency_files',
      placeholder: 'Exemple : package.json,yarn.lock (séparés par une virgule)',
      agreement: 'plural',
      defaultValue: project.dependency_files ? project.dependency_files.join(',') : '',
    },
    {
      title: "Script d'installation",
      name: 'install_script',
      defaultValue: project.install_script || '',
    },
    {
      title: 'Script de build',
      name: 'build_script',
      defaultValue: project.build_script || '',
    },
    {
      title: 'Script de run',
      name: 'start_script',
      defaultValue: project.start_script || '',
    },
    {
      title: 'Dossier statique',
      placeholder: 'Exemple images',
      name: 'static_folder',
      defaultValue: project.static_folder || '',
    },
    {
      type: 'number',
      title: "Port d'écoute",
      name: 'exposed_port',
      defaultValue: project.exposed_port || '',
      min: 80,
      max: 65535,
    },
  ].filter(question => !ignoreField[question.name]);

  const steps = [
    {
      name: 'Url',
      onSubmit: onSubmit({
        repository_url: Boolean,
      }),
      content: <Repositories onSelect={repository => setRepo(repository)} />,
    },
    {
      name: 'Langage',
      onSubmit: onSubmit(
        {
          name: name => checkProject(name).then(({ data }) => data === false),
          branch: Boolean,
          template_id: Boolean,
          external_domain_names: edn => typeof edn === 'string' && !edn.includes('/'),
        },
        p => ({ ...p, ...template, external_domain_names: split(p.external_domain_names) }),
      ),
      content: (
        <>
          <Input
            title="Nom du projet"
            name="name"
            error={error === 'name'}
            errorMessage="Ce projet existe déjà, trouve un autre nom 😉"
            defaultValue={repo ? repo.name : ''}
            required
          />
          <Select
            title="Sur quelle branche veux-tu que ton projet soit déployé ?"
            name="branch"
            error={error === 'branch'}
            required
            data={branches.map(({ name }) => ({ value: name }))}
          />
          <Input
            title="Autres noms de domaine"
            name="external_domain_names"
            placeholder="Exemple : mon-site.fr,www.mon-site.fr (séparés par une virgule, sans https ni / à la fin)"
            agreement="plural"
            defaultValue={project.external_domain_names ? project.external_domain_names.join(',') : ''}
          />
          <TemplatePicker onSelect={setTemplate} />
        </>
      ),
    },
    questions.length && {
      name: 'Variables',
      onSubmit: onSubmit(
        {
          exposed_port: value => !value || !Number.isNaN(Number(value)),
        },
        ({ dependency_files, system_dependencies, use_https, exposed_port }) => ({
          dependency_files: split(dependency_files),
          system_dependencies: split(system_dependencies),
          use_https: use_https === 'https',
          exposed_port: Number(exposed_port) || null,
        }),
      ),
      content: <Variables project={project} questions={questions} />,
    },
    ignoreField.databases
      ? null
      : {
          name: 'Bdd',
          onSubmit: onSubmit({}, ({ databases }) => ({
            databases: !databases || Array.isArray(databases) ? [] : Object.values(databases),
          })),
          content: <Database />,
        },
    {
      name: 'Environnement',
      onSubmit: onSubmit({}, null, () =>
        postProject(project)
          .then(({ data: { id } }) => activateProject(id))
          .then(() => setLoading(false)),
      ),
      content: <Envs />,
    },
    {
      name: 'Envoyer',
      content: loading ? 'Ton site va être déployé. Vérifie dans quelques instants' : 'Ton site est déployé',
    },
  ].filter(Boolean);

  return (
    <Steps step={step}>
      {steps.map((s, index) => (
        <Box
          key={s.name}
          className={block}
          {...(s.onSubmit
            ? {
                as: 'form',
                onSubmit: s.onSubmit,
              }
            : {})}
        >
          {s.content}
          {index !== steps.length - 1 && (
            <div className={direction}>
              {index === 0 ? (
                <div />
              ) : (
                <ArrowButton type="button" onClick={previousStep} direction="previous">
                  {steps[index - 1].name}
                </ArrowButton>
              )}
              <ArrowButton type="submit">{steps[index + 1].name}</ArrowButton>
            </div>
          )}
        </Box>
      ))}
    </Steps>
  );
};

export default New;
