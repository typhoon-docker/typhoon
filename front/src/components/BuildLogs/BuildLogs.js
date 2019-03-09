import React, { useRef } from 'react';

import hljs from 'highlight.js/lib/highlight';
import accesslog from 'highlight.js/lib/languages/accesslog';

import { getBuildLogs } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';
import cx from '/utils/className';

import { lines, line } from './BuildLogs.css';

hljs.registerLanguage('accesslog', accesslog);

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
            <pre ref={ref} className={cx(lines, 'accesslog')}>
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
