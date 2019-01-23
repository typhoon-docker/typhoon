import React from 'react';

import { content } from './EmptyProject.css';

const EmptyProject = () => (
  <div className={content}>
    <p>{"Oh non, tu n'as pas encore créé de projets 😥"}</p>
    <p>{'Rejoins nous vite 😎'}</p>
  </div>
);

export default EmptyProject;
