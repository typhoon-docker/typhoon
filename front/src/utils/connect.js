import { decode, sign } from 'jsonwebtoken';
import { shouldMock } from './env';

const storage = sessionStorage;

const TOKEN_KEY = 'token-hjqbgk-oiqjwe-1-4.0';

export const saveToken = token => storage.setItem(TOKEN_KEY, token);

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

export const isConnected = () => {
  try {
    const token = getToken();
    if (token.exp < Date.now() / 1000 || !token.exp) {
      throw new Error('');
    }
    return true;
  } catch (error) {
    return false;
  }
};

if (shouldMock) {
  import('./mock/user.json').then(mockUser => {
    const iat = parseInt(Date.now() / 1000, 10) - 60;
    const exp = iat + 3600;
    saveToken(
      sign(
        {
          user: mockUser,
          iat,
          exp,
        },
        'secret',
      ),
    );
  });
}
