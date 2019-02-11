export const formDataToArray = formData => Array.from(formData.entries());

const regex = /\./g;

const apply = (obj, key, value) => {
  if (!key.match(regex)) {
    return { ...obj, [key]: value };
  }

  let mem = obj;
  const path = key.split('.');
  const [lastKey] = path.splice(path.length - 1);
  path.forEach(p => {
    if (!(p in mem)) {
      mem[p] = {};
    }
    mem = mem[p];
  });
  mem[lastKey] = value;
  return obj;
};

export const arrayToJSON = array => array.reduce((acc, [key, value]) => apply(acc, key, value), {});

export const formDataToJSON = formData => arrayToJSON(formDataToArray(formData));
