import React from 'react';
import { Redirect } from '@reach/router';
import { saveToken } from '/utils/connect';

const CallbackViaRezo = () => {
  const qs = new URLSearchParams(window.location.search);
  const token = qs.get('token');
  if (token !== null) {
    saveToken(token);
  }
  return <Redirect to="/" noThrow />;
};

export default CallbackViaRezo;
