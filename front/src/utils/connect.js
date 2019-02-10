import { useEffect, useState } from 'react';
import { decode, sign } from 'jsonwebtoken';
import { createCup } from 'manatea';

import { shouldMock } from './env';

const TOKEN_KEY = 'token-hjqbgk-oiqjwe-1-4.0';
const storage = sessionStorage;

export const tokenCup = createCup(null);

tokenCup.on(token => storage.setItem(TOKEN_KEY, token));

export const saveToken = token => tokenCup(token);

export const removeToken = () => storage.removeItem(TOKEN_KEY);

export const getRawToken = () => storage.getItem(TOKEN_KEY);

export const getToken = () => decode(getRawToken());

export const hasToken = () => {
  try {
    getToken();
    return true;
  } catch (error) {
    return false;
  }
};

export const isConnected = (token = getToken()) => {
  try {
    if (!token.exp || token.exp < Date.now() / 1000) {
      throw new Error('');
    }
    return true;
  } catch (error) {
    return false;
  }
};

export const isAdmin = (token = getToken()) => {
  try {
    if (!token.user.scope.includes('admin')) {
      throw new Error('');
    }
    return true;
  } catch (error) {
    return false;
  }
};

// React Hooks
const hook = fn => () => {
  const [value, setValue] = useState(fn());
  useEffect(() => {
    const listener = tokenCup.on(token => {
      setValue(fn(decode(token)));
    });
    setValue(fn());
    return listener;
  }, []);
  return value;
};

export const useIsConnected = hook(isConnected);
export const useIsAdmin = hook(isAdmin);

if (shouldMock) {
  import('./mock/user.json').then(mockUser => {
    const iat = parseInt(Date.now() / 1000, 10) - 60;
    const exp = iat + 3600;
    saveToken(
      sign(
        {
          user: mockUser,
          exp,
        },
        'secret',
      ),
    );
  });
}
