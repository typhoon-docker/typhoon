import React from 'react';
import { Redirect } from '@reach/router';
import { newProjectCup } from '/utils/project';

const CallbackGitHub = () => {
  const qs = new URLSearchParams(window.location.search);
  const token = qs.get('token');
  if (token === null) {
    return <Redirect to="/error" noThrow />;
  }
  newProjectCup(project => ({ ...project, repository_token: token }));
  return <Redirect to="/new" noThrow />;
};

export default CallbackGitHub;
