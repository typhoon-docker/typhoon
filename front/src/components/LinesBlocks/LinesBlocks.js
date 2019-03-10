import React from 'react';

import { getBuildLogs, getDockerFiles } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';

import { lbh2, lbh3, block, lines, line } from './LinesBlocks.css';

const OneLinesBlock = (keyName, linesText) => {
  if (!linesText || keyName.startsWith('_')) {
    return null;
  }
  return (
    <div key={keyName}>
      <h3 className={lbh3}>{keyName}</h3>
      <div className={block}>
        <pre className={lines}>
          {linesText
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
    </div>
  );
};

const LinesBlocks = (linesData, title) => {
  return (
    <>
      <h2 className={lbh2}>{title}</h2>
      {Object.keys(linesData).map(keyName => OneLinesBlock(keyName, linesData[keyName]))}
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
