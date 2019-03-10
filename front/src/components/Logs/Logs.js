import React, { useEffect, useRef } from 'react';

import hljs from 'highlight.js/lib/highlight';
import accesslog from 'highlight.js/lib/languages/accesslog';

import { getLogs } from '/utils/typhoonAPI';
import useAxios from '/utils/useAxios';
import cx from '/utils/className';

import { button, lines, line } from './Logs.css';

hljs.registerLanguage('accesslog', accesslog);

const Logs = ({ projectID }) => {
  const [logs, , refetch] = useAxios(getLogs(projectID), '', [projectID]);
  const ref = useRef(null);

  useEffect(() => {
    if (logs) {
      Array.from(ref.current.children).forEach(el => {
        hljs.highlightBlock(el);
      });
    }
  }, [logs]);

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
      <pre ref={ref} className={cx(lines, 'accesslog')}>
        <code>
          {logs
            .trim()
            .split('\n')
            .map(log => (
              <span className={line} key={log}>
                {log}
              </span>
            ))}
        </code>
      </pre>
    </>
  );
};

export default Logs;
