import h from '/utils/h';

import { Link } from '@reach/router';

import ArrowButton from '/components/ArrowButton';
import StatusBubbles from '/components/StatusBubbles';

import { details, summary, content } from './Project.css';
import { highlightable } from '/styles/highlightable.css';

import cx from '/utils/className/';
import { useProperty } from '/utils/hooks';

const Project = ({ project, onSelect, selected }) => {
  const { id, name, repository_url, template_id } = project;
  const el = useProperty(
    {
      '--border-color': `rgb(var(--${template_id}))`,
    },
    [template_id],
  );

  return (
    <details key={id} open={selected} className={details}>
      <summary onClick={onSelect} className={cx(summary, highlightable)}>
        {name}
      </summary>
      <div className={content} ref={el}>
        <span style={{ gridArea: 'name' }}>
          <a href={`https://${name.toLowerCase()}.typhoon.viarezo.fr/`} target="_blank" rel="noopener noreferrer">
            {`https://${name.toLowerCase()}.typhoon.viarezo.fr/`}
          </a>{' '}
          <a href={repository_url} target="_blank" rel="noopener noreferrer">
            (source code)
          </a>
        </span>
        <span style={{ gridArea: 'button' }}>
          <ArrowButton
            as={Link}
            to={`/project/${id}`}
            color="text"
            style={{ fontSize: '1.2em', textDecoration: 'none' }}
          >
            Voir le projet
          </ArrowButton>
        </span>
        <StatusBubbles projectID={project.id} style={{ gridArea: 'bubbles' }} />
      </div>
    </details>
  );
};

export default Project;
