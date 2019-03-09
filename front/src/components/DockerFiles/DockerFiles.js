import React, { useRef } from 'react';

import { getDockerFiles } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';

import { lines, line } from './DockerFiles.css';

const DockerFiles = ({ projectID }) => {
  const dockerFilesData = useAxios(getDockerFiles(projectID), { logs: '' }, [projectID])[0];
  const ref = useRef(null);

  return (
    <>
      <h2>Docker Files</h2>
      {Object.keys(dockerFilesData).map(keyName => {
        if (!dockerFilesData[keyName] || keyName === 'project') {
          return null;
        }
        return (
          <div key={keyName}>
            <h3>{keyName}</h3>
            <pre ref={ref} className={lines}>
              <code>
                {dockerFilesData[keyName]
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

export default DockerFiles;
