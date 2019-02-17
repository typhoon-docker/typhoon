import React, { useState, useEffect } from 'react';

import { statusProject } from '/utils/typhoonAPI';

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
  const [containers, setContainers] = useState([]);

  useEffect(() => {
    statusProject(projectID).then(({ data }) => setContainers(data));
  }, [projectID]);

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
