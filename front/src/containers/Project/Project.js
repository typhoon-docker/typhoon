import React from 'react';

import ArrowButton from '/components/ArrowButton';
import StatusBubbles from '/components/StatusBubbles';

import { details, summary, content } from './Project.css';
import { highlightable } from '/styles/highlightable.css';

import cx from '/utils/className';
import useProperty from '/utils/useProperty';

const Project = ({ project, onSelect, selected }) => {
  const { id, name, repository_url, template_id } = project;
  const el = useProperty(
    () => ({
      '--border-color': `rgb(var(--${template_id}))`,
    }),
    [template_id],
  );

  return (
    <details key={id} open={selected} className={details}>
      <summary onClick={onSelect} className={cx(summary, highlightable)}>
        {name}
      </summary>
      <div className={content} ref={el}>
        <div>
          <a href={`https://${name}.typhoon.viarezo.fr/`} target="_blank" rel="noopener noreferrer">
            {`https://${name}.typhoon.viarezo.fr/`}
          </a>{' '}
          <a href={repository_url} target="_blank" rel="noopener noreferrer">
            (source code)
          </a>
          <ArrowButton>Voir le projet</ArrowButton>
        </div>
        <StatusBubbles projectID={project.id} />
      </div>
    </details>
  );
};

export default Project;
