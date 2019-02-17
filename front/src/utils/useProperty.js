import { useRef, useEffect } from 'react';

const useProperty = (cb, mem) => {
  const ref = useRef(null);

  useEffect(() => {
    Object.entries(cb()).forEach(([key, value]) => {
      ref.current.style.setProperty(key, value);
    });
  }, mem);

  return ref;
};

export default useProperty;
