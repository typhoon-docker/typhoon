import React from 'react';

import { statusProject } from '/utils/typhoonAPI';
import { useAxios } from '/utils/hooks';

import { bubble } from './StatusBubbles.css';

const colors = {
  created: 'dodgerblue',
  restarting: 'lightskyblue',
  running: 'yellowgreen',
  removing: 'tomato',
  paused: 'grey',
  exited: 'gainsboro',
  dead: 'black',
};

const StatusBubbles = ({ projectID, ...props }) => {
  const [containers] = useAxios(statusProject(projectID), [], [projectID]);

  return (
    <div {...props}>
      {containers.map(({ id, state, name }) => (
        <div
          key={id}
          title={`${name} : ${state}`}
          className={bubble}
          style={{
            backgroundColor: colors[state],
          }}
        />
      ))}
    </div>
  );
};

export default StatusBubbles;
