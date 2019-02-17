import React from 'react';
import { Redirect } from '@reach/router';
import { saveToken, retrieveLocation } from '/utils/connect';

const CallbackViaRezo = () => {
  const qs = new URLSearchParams(window.location.search);
  const token = qs.get('token');
  if (token !== null) {
    saveToken(token);
  }
  const location = retrieveLocation();

  return <Redirect to={location || '/'} noThrow />;
};

export default CallbackViaRezo;
