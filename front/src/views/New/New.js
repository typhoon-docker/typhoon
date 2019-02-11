import React, { useState, useEffect } from 'react';

import Repositories from '/containers/Repositories/';

import Steps from '/components/Steps/';
import Box from '/components/Box/';
import { Text, Select } from '/components/Input/';
import ArrowButton from '/components/ArrowButton/';
import TemplatePicker from '/components/TemplatePicker/';

import { newProjectCup } from '/utils/project';
import { formDataToArray, arrayToJSON } from '/utils/formData';
import { checkProject } from '/utils/typhoonAPI';
import { getBranches } from '/utils/githubAPI';

import { block, next } from './New.css';

const New = () => {
  const [step, setStep] = useState(0);
  const [error, setError] = useState(null);
  const [repo, setRepo] = useState(null);
  const [template, setTemplate] = useState({});
  const [branches, setBranches] = useState([]);

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

  const onSubmit = (verifies, cb = () => {}) => event => {
    event.preventDefault();
    const newDataArray = formDataToArray(new FormData(event.target));
    if (Object.keys(verifies).length > newDataArray.length) {
      return;
    }
    if (Object.keys(verifies).length < newDataArray.length) {
      console.warn("Many fields aren't verified");
    }
    Promise.all(newDataArray.map(([key, value]) => verifies[key](value)))
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
        const newData = arrayToJSON(array);
        newProjectCup(project => ({ ...project, ...newData }));
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
        <div className={next}>
          <ArrowButton type="submit">Continuer</ArrowButton>
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
          () => newProjectCup(project => ({ ...project, ...template })),
        )}
      >
        <Text
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

        <div className={next}>
          <ArrowButton type="submit">Continuer</ArrowButton>
        </div>
      </Box>
      <Box as="form" className={block}>
        <div className={next}>
          <ArrowButton type="submit">Continuer</ArrowButton>
        </div>
      </Box>
    </Steps>
  );
};

export default New;
