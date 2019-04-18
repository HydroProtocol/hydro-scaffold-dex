import BigNumber from 'bignumber.js';

export const sleep = time => new Promise(r => setTimeout(r, time));

export const toUnitAmount = (amount, decimals) => {
  return new BigNumber(amount).div(Math.pow(10, decimals));
};

export const isTokenApproved = allowance => {
  return allowance.gt(10 ** 30);
};
