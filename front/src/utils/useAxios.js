import { useState, useEffect } from 'react';

const useAxios = (query, defaultValue, mem) => {
  const [value, setValue] = useState(defaultValue);

  useEffect(() => {
    query.then(({ data }) => setValue(data));
  }, mem);

  return [value, setValue, () => query.then(({ data }) => setValue(data))];
};

export default useAxios;
