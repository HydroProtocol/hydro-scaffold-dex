import { setConfigs } from './config';
import { getTokenBalance } from '../lib/wallet';
import env from '../lib/env';

export let hotDiscountRules = [];

export const loadHotDiscountRules = async () => {
  hotDiscountRules = [[5000, 1], [20000, 0.9], [100000, 0.8], [500000, 0.7], [2000000, 0.6], [-1, 0.5]];
};

export const getHotTokenAmount = () => {
  return async (dispatch, getState) => {
    const hotContract = env.HYDRO_TOKEN_ADDRESS;
    if (!hotContract) {
      return;
    }

    const address = getState().account.get('address');
    if (!address) {
      return;
    }
    const hotTokenAmount = await getTokenBalance(hotContract, address);
    dispatch(setConfigs({ hotTokenAmount }));
  };
};
