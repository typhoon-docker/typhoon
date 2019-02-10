const formDataToJSON = formData =>
  Array.from(formData.entries()).reduce(
    (memo, pair) => ({
      ...memo,
      [pair[0]]: (console.log(pair), pair[1]),
    }),
    {},
  );

export default formDataToJSON;
