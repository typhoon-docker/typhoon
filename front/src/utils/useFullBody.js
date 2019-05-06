import { useEffect } from 'react';
import { full } from '../App.css';

const useFullBody = () => {
  useEffect(() => {
    document.body.classList.add(full);
    return () => {
      document.body.classList.remove(full);
    };
  }, []);
};

export default useFullBody;
