import React, { useState } from 'react';

import Steps from '/components/Steps/';
import Box from '/components/Box/';
import Input from '/components/Input/';
import ArrowButton from '/components/ArrowButton/';
import Repositories from '/components/Repositories/';

import { newProjectCup } from '/utils/project';
import { formDataToArray, arrayToJSON } from '/utils/formData';
import { checkProject } from '/utils/typhoonAPI';

const New = () => {
  const [step, setStep] = useState(0);
  const [error, setError] = useState(null);
  const [repo, setRepo] = useState(null);
  // const previousStep = () => setStep(step - 1);
  const nextStep = () => setStep(step + 1);

  const onSubmit = verifies => event => {
    event.preventDefault();
    const newDataArray = formDataToArray(new FormData(event.target));
    if (Object.keys(verifies).length !== newDataArray.length) {
      return;
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
      .catch(console.warn);
  };

  return (
    <Steps step={step}>
      <Box
        as="form"
        onSubmit={onSubmit({
          repo: value => {
            try {
              const repository = JSON.parse(value);
              setRepo(repository);
              return true;
            } catch (e) {
              return false;
            }
          },
        })}
      >
        <Repositories />
      </Box>
      <Box
        as="form"
        onSubmit={onSubmit({
          name: name => checkProject(name).then(({ data }) => data === false),
        })}
      >
        <Input
          title="Nom du projet"
          name="name"
          error={error === 'name'}
          errorMessage="Ce projet existe déjà, trouvez un autre nom"
          defaultValue={repo ? repo.full_name : ''}
          required
        />
        <ArrowButton type="submit">Next</ArrowButton>
      </Box>
      <Box as="form">test</Box>
    </Steps>
  );
};

export default New;
