import React, { useEffect, useRef } from 'react';

import { getLogs } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';

import { button, lines } from './Logs.css';

const Logs = ({ projectID }) => {
  const [logs, , refetch] = useAxios(getLogs(projectID), '', [projectID]);
  const ref = useRef(null);

  useEffect(() => {
    ref.current.scrollTop = ref.current.scrollHeight;
  }, [logs]);

  return (
    <>
      <h2>
        Logs{' '}
        <button type="button" onClick={refetch} className={button}>
          <span role="img" aria-label="retry">
            🔄
          </span>
        </button>
      </h2>
      <pre ref={ref} className={lines}>
        {logs}
      </pre>
    </>
  );
};

export default Logs;
