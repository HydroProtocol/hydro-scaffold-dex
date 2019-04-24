import { Map } from 'immutable';
import BigNumber from 'bignumber.js';

export const stateUtils = {
  getTokensInfo: (state, accountAddress) => {
    return state.account.getIn(['tokensInfo', accountAddress], Map());
  },
  getTokenInfo: (state, accountAddress, tokenSymbol) => {
    return state.account.getIn(['tokensInfo', accountAddress, tokenSymbol], Map());
  },
  getTokenBalance: (state, accountAddress, tokenSymbol) => {
    return state.account.getIn(['tokensInfo', accountAddress, tokenSymbol, 'balance'], new BigNumber('0'));
  },
  getTokenLockedBalance: (state, accountAddress, tokenSymbol) => {
    return state.account.getIn(['tokensInfo', accountAddress, tokenSymbol, 'lockedBalance'], new BigNumber('0'));
  },
  getTokenAvailableBalance: (state, accountAddress, tokenSymbol) => {
    const balance = stateUtils.getTokenBalance(state, accountAddress, tokenSymbol);
    const lockedBalance = stateUtils.getTokenLockedBalance(state, accountAddress, tokenSymbol);
    return balance.minus(lockedBalance);
  }
};
