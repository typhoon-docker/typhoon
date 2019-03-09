import React from 'react';

import { getBuildLogs, getDockerFiles } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';

import { block, lines, line } from './LinesBlocks.css';

const LinesBlocks = (linesData, title) => {
  return (
    <>
      <h2>{title}</h2>
      {Object.keys(linesData).map(keyName => {
        if (!linesData[keyName] || keyName.startsWith('_')) {
          return null;
        }
        return (
          <div key={keyName}>
            <h3>{keyName}</h3>
            <div className={block}>
              <pre className={lines}>
                {linesData[keyName]
                  .trim()
                  .split('\n')
                  .map((lineText, i) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <span className={line} key={i}>
                      {lineText}
                    </span>
                  ))}
              </pre>
            </div>
          </div>
        );
      })}
    </>
  );
};

const BuildLogs = ({ projectID }) => {
  const buildLogsData = useAxios(getBuildLogs(projectID), {}, [projectID])[0];
  return LinesBlocks(buildLogsData, 'Build Logs');
};

const DockerFiles = ({ projectID }) => {
  const dockerFilesData = useAxios(getDockerFiles(projectID), {}, [projectID])[0];
  return LinesBlocks(dockerFilesData, 'Docker Files');
};

export { BuildLogs, DockerFiles };
