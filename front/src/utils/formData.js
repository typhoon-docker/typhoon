export const formDataToArray = formData => Array.from(formData.entries());

export const arrayToJSON = array =>
  array.reduce(
    (acc, [key, value]) => ({
      ...acc,
      [key]: value,
    }),
    {},
  );

export const formDataToJSON = formData => arrayToJSON(formDataToArray(formData));
