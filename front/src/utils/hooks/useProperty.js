import { useRef, useEffect } from 'react';

const useProperty = (object, mem) => {
  const ref = useRef(null);

  useEffect(() => {
    Object.entries(object).forEach(([key, value]) => {
      ref.current.style.setProperty(key, value);
    });
  }, mem);

  return ref;
};

export default useProperty;
