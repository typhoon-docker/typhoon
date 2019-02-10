import React, { useState, useRef } from 'react';

import Steps from '/components/Steps/';
import Box from '/components/Box/';

import { newProjectCup } from '/utils/project';
import formDataToJSON from '/utils/formDataToJSON';
import { checkProject } from '/utils/axios';

const New = () => {
  const [step, setStep] = useState(0);
  // const previousStep = () => setStep(step - 1);
  const nextStep = () => setStep(step + 1);

  const nameRef = useRef(null);
  const onSubmit = event => {
    event.preventDefault();
    newProjectCup(project => ({ ...project, ...formDataToJSON(new FormData(event.target)) }));
    nextStep();
  };

  return (
    <Steps step={step}>
      <Box as="form" onSubmit={onSubmit}>
        <input
          type="text"
          name="name"
          required
          ref={nameRef}
          verify={name => checkProject(name).then(data => data === 'true')}
        />
        <button type="submit">Next</button>
      </Box>
      <Box as="form">test</Box>
    </Steps>
  );
};

export default New;
