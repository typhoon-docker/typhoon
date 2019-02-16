import React, { useState, useEffect } from 'react';

import Repositories from '/containers/Repositories/';
import Variables from '/containers/Variables/';
import Database from '/containers/Database/';
import Envs from '/containers/Envs/';

import Steps from '/components/Steps/';
import Box from '/components/Box/';
import { Input, Select } from '/components/Input/';
import ArrowButton from '/components/ArrowButton/';
import TemplatePicker from '/components/TemplatePicker/';

import { useProject } from '/utils/project';
import { formDataToArray, arrayToJSON } from '/utils/formData';
import { checkProject, postProject, activateProject } from '/utils/typhoonAPI';
import { getBranches } from '/utils/githubAPI';

import { block, direction } from './New.css';

const New = () => {
  const [step, setStep] = useState(0);
  const [error, setError] = useState(null);
  const [repo, setRepo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [template, setTemplate] = useState({});
  const [branches, setBranches] = useState([]);
  const [project, setProject] = useProject();

  useEffect(() => {
    if (repo) {
      getBranches(repo)
        .then(({ data }) => {
          setBranches(data);
        })
        .catch(console.warn);
    }
  }, [repo]);

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
        let newData = arrayToJSON(array);
        newData = { ...newData, ...(transform || (x => x))(newData) };
        setProject(p => ({ ...p, ...newData }));
        nextStep();
      })
      .then(cb)
      .catch(console.warn);
  };

  return (
    <Steps step={step}>
      <Box
        as="form"
        className={block}
        onSubmit={onSubmit({
          repository_url: Boolean,
        })}
      >
        <Repositories onSelect={repository => setRepo(repository)} />
        <div className={direction}>
          <div />
          <ArrowButton type="submit">Langage</ArrowButton>
        </div>
      </Box>
      <Box
        as="form"
        className={block}
        onSubmit={onSubmit(
          {
            name: name => checkProject(name).then(({ data }) => data === false),
            branch: Boolean,
            template_id: Boolean,
          },
          p => ({ ...p, ...template }),
        )}
      >
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
        <TemplatePicker onSelect={setTemplate} />

        <div className={direction}>
          <ArrowButton type="button" onClick={previousStep} direction="previous">
            Url
          </ArrowButton>
          <ArrowButton type="submit">Variables</ArrowButton>
        </div>
      </Box>
      <Box
        as="form"
        className={block}
        onSubmit={onSubmit(
          {
            exposed_port: value => !value || !Number.isNaN(Number(value)),
          },
          ({ external_domain_names, dependency_files, system_dependencies, use_https, exposed_port }) => ({
            external_domain_names: (external_domain_names || '').split(','),
            dependency_files: (dependency_files || '').split(','),
            system_dependencies: (system_dependencies || '').split(','),
            use_https: use_https === 'https',
            exposed_port: Number(exposed_port) || null,
          }),
        )}
      >
        <Variables project={project} />
        <div className={direction}>
          <ArrowButton type="button" onClick={previousStep} direction="previous">
            Langage
          </ArrowButton>
          <ArrowButton type="submit">Bdd</ArrowButton>
        </div>
      </Box>
      <Box
        as="form"
        onSubmit={onSubmit({}, ({ databases }) => ({
          databases: !databases || Array.isArray(databases) ? [] : Object.values(databases),
        }))}
      >
        <Database />
        <div className={direction}>
          <ArrowButton type="button" onClick={previousStep} direction="previous">
            Url
          </ArrowButton>
          <ArrowButton type="submit">Variables</ArrowButton>
        </div>
      </Box>
      <Box
        as="form"
        onSubmit={onSubmit({}, null, () =>
          postProject(project)
            .then(({ data: { id } }) => activateProject(id))
            .then(() => setLoading(false)),
        )}
      >
        <Envs />
        <div className={direction}>
          <ArrowButton type="button" onClick={previousStep} direction="previous">
            Db
          </ArrowButton>
          <ArrowButton type="submit">Valider</ArrowButton>
        </div>
      </Box>
      <Box>{loading ? 'Ton site va être déployé. Vérifie dans quelques instants' : 'Ton site est déployé'}</Box>
    </Steps>
  );
};

export default New;
