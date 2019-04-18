import { Map } from 'immutable';
import BigNumber from 'bignumber.js';

const initialState = Map({
  WETH: {
    address: '0x4a817489643A89a1428b2DD441c3fbe4DBf44789',
    symbol: 'WETH',
    decimals: 18
  },
  hotTokenAmount: new BigNumber(0),
  websocketConnected: false,
  web3NetworkID: null
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
