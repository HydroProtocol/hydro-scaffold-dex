import axios from 'axios';
import env from './env';
import { store } from '../index';
import { cleanLoginDate, loadAccountHydroAuthentication } from './session';
import { logout } from '../actions/account';
import { getSelectedAccount } from '@gongddex/hydro-sdk-wallet';

const getAxiosInstance = () => {
  const state = store.getState();
  const selectedAccount = getSelectedAccount(state);
  const address = selectedAccount ? selectedAccount.get('address') : null;
  const hydroAuthentication = loadAccountHydroAuthentication(address);
  let instance;

  if (hydroAuthentication) {
    instance = axios.create({
      headers: {
        'Hydro-Authentication': hydroAuthentication
      }
    });
  } else {
    instance = axios;
  }

  instance.interceptors.response.use(function(response) {
    if (response.data && response.data.status === -11) {
      if (address) {
        store.dispatch(logout(address));
        cleanLoginDate(address);
      }
    }
    return response;
  });

  return instance;
};

const _request = (method, url, ...args) => {
  return getAxiosInstance()[method](`${env.API_ADDRESS}${url}`, ...args);
};

const _coinBaseRequest = (method, url, ...args) => {
  const instance = axios.create({
    baseURL: env.COIN_BASE_API_ADDRESS,
    headers: {
      'X-CoinAPI-Key': env.COIN_BASE_API_KEY
    }
  });

  return instance[method](url, ...args);
}

const api = {
  get: (url, ...args) => _request('get', url, ...args),
  delete: (url, ...args) => _request('delete', url, ...args),
  head: (url, ...args) => _request('head', url, ...args),
  post: (url, ...args) => _request('post', url, ...args),
  put: (url, ...args) => _request('put', url, ...args),
  patch: (url, ...args) => _request('patch', url, ...args),
  coinBaseGet: (url, ...args) => _coinBaseRequest('get', url, ...args)
};

export default api;
