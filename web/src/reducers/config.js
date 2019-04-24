import { Map } from 'immutable';
import BigNumber from 'bignumber.js';
import env from '../lib/env';

const initialState = Map({
  WETH: {
    address: env.WETH_TOKEN_ADDRESS,
    symbol: 'WETH',
    decimals: 18
  },
  hotTokenAmount: new BigNumber(0),
  websocketConnected: false
});

export default (state = initialState, action) => {
  switch (action.type) {
    case 'SET_CONFIGS':
      for (const key of Object.keys(action.payload)) {
        state = state.set(key, action.payload[key]);
      }
      return state;
    default:
      return state;
  }
};
