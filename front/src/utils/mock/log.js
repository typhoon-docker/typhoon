const analyse = (group, [status, body, headers]) => {
  console.group(group);
  console.log('Status:', status);
  if (body) {
    if (Array.isArray(body)) {
      console.log('Body:');
      console.table(body);
    } else {
      console.log('Body', body);
    }
  }
  if (headers) {
    console.log('Headers:', headers);
  }
  console.groupEnd(group);
};

const log = (...args) => async config => {
  const { method, baseURL, url } = config;
  const group = `Mocked ${method.toUpperCase()} - ${url.substring(baseURL.length)}`;
  if (args.length === 1 && typeof args[0] === 'function') {
    const response = await args[0](config);
    analyse(group, response);
    return response;
  }
  analyse(group, args);
  return args;
};

export default log;
