import React from 'react';

import { getBuildLogs, getDockerFiles } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';

import { lbh2, lbh3, block, lines, line } from './LinesBlocks.css';

const TextBlock = ({ title, text }) => {
  if (!text || title.startsWith('_')) {
    return null;
  }
  return (
    <>
      <h3 className={lbh3}>{title}</h3>
      <div className={block}>
        <pre className={lines}>
          {text
            .trim()
            .split('\n')
            .map((lineText, i) => (
              // The key is the array index because it will not be refreshed, no need for tracking
              // eslint-disable-next-line react/no-array-index-key
              <span className={line} key={i}>
                {lineText}
              </span>
            ))}
        </pre>
      </div>
    </>
  );
};

const textMultipleBlocks = (linesData, title) => {
  return (
    <>
      <h2 className={lbh2}>{title}</h2>
      {Object.keys(linesData).map(keyName => (
        <TextBlock key={keyName} title={keyName} text={linesData[keyName]} />
      ))}
    </>
  );
};

const BuildLogs = ({ projectID }) => {
  const buildLogsData = useAxios(getBuildLogs(projectID), {}, [projectID])[0];
  return textMultipleBlocks(buildLogsData, 'Build Logs');
};

const DockerFiles = ({ projectID }) => {
  const dockerFilesData = useAxios(getDockerFiles(projectID), {}, [projectID])[0];
  return textMultipleBlocks(dockerFilesData, 'Docker Files');
};

export { BuildLogs, DockerFiles };
