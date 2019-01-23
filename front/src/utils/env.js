const MOCK = process.env.MOCK ? process.env.MOCK.toLowerCase() : '';
const NODE_ENV = process.env.NODE_ENV ? process.env.NODE_ENV.toLowerCase() : '';

export const shouldMock = MOCK === 'true';
export const isProd = NODE_ENV === 'production' || NODE_ENV === 'prod';
export const isDev = NODE_ENV === 'development' || NODE_ENV === 'dev' || NODE_ENV === '';
export const isTest = NODE_ENV === 'test';
