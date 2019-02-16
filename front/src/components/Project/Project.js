import React, { useRef, useEffect } from 'react';

import { details, summary, content } from './Project.css';
import { highlightable } from '/styles/highlightable.css';

import cx from '/utils/className';

const Project = ({ project, onSelect, selected }) => {
  const { id, name, repository_url, template_id } = project;
  const el = useRef(null);

  useEffect(() => {
    el.current.style.setProperty('--border-color', `rgb(var(--${template_id}))`);
  });

  return (
    <details key={id} open={selected} className={details}>
      <summary onClick={onSelect} className={cx(summary, highlightable)}>
        {name}
      </summary>
      <div className={content} ref={el}>
        Projet en {template_id}, accessible via{' '}
        <a href={`https://${name}.typhoon.viarezo.fr/`} target="_blank" rel="noopener noreferrer">
          {`https://${name}.typhoon.viarezo.fr/`}
        </a>{' '}
        <a href={repository_url} target="_blank" rel="noopener noreferrer">
          (source code)
        </a>
      </div>
    </details>
  );
};

export default Project;
