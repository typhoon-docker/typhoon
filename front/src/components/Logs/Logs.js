import h from '/utils/h';

import { useEffect, useRef } from 'react';

import hljs from 'highlight.js/lib/highlight';
import accesslog from 'highlight.js/lib/languages/accesslog';

import { getLogs } from '/utils/typhoonAPI';
import { useAxios } from '/utils/hooks';
import cx from '/utils/className/';

import { button, lines, line } from './Logs.css';

hljs.registerLanguage('accesslog', accesslog);

const createRawLogs = logs => {
  const code = document.createElement('code');
  logs
    .trim()
    .split('\n')
    .forEach(log => {
      const span = document.createElement('span');
      span.classList.add(line);
      span.textContent = log;
      code.appendChild(span);
    });
  return { __html: code.outerHTML };
};

const Logs = ({ projectID }) => {
  const [logs, , refetch] = useAxios(() => getLogs(projectID), '', [projectID]);
  const logRef = useRef(null);

  useEffect(() => {
    if (logs.length && logRef.current.children) {
      Array.from(logRef.current.children).forEach(el => {
        hljs.highlightBlock(el);
      });
    }
  }, [logs]);

  useEffect(() => {
    logRef.current.scrollTop = logRef.current.scrollHeight;
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
      <pre ref={logRef} className={cx(lines, 'accesslog')} dangerouslySetInnerHTML={createRawLogs(logs)} />
    </>
  );
};

export default Logs;
