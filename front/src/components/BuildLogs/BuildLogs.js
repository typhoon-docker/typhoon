import React, { useRef } from 'react';

import { getBuildLogs } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';

import { lines, line } from './BuildLogs.css';

const BuildLogs = ({ projectID }) => {
  const logsData = useAxios(getBuildLogs(projectID), { logs: '' }, [projectID])[0].logs;
  const ref = useRef(null);

  return (
    <>
      <h2>Build Logs</h2>
      {Object.keys(logsData).map(keyName => {
        if (!logsData[keyName]) {
          return null;
        }
        return (
          <div key={keyName}>
            <h3>{keyName}</h3>
            <pre ref={ref} className={lines}>
              <code>
                {logsData[keyName]
                  .trim()
                  .split('\n')
                  .map((log, j) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <span className={line} key={j}>
                      {log}
                    </span>
                  ))}
              </code>
            </pre>
          </div>
        );
      })}
    </>
  );
};

export default BuildLogs;
