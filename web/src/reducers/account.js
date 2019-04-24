import { Map, OrderedMap } from 'immutable';
import { BigNumber } from 'bignumber.js';

export const initState = Map({
  isLoggedIn: Map(),
  tokensInfo: Map(),
  approving: Map(),
  orders: OrderedMap(),
  trades: OrderedMap(),
  transactions: OrderedMap()
});

const initialTokenInfo = Map({
  balance: new BigNumber('0'),
  allowance: new BigNumber('0'),
  address: '',
  decimals: 0,
  lockedBalance: new BigNumber('0')
});

export default (state = initState, action) => {
  switch (action.type) {
    case 'UPDATE_TOKEN_LOCKED_BALANCES': {
      const { lockedBalances, accountAddress } = action.payload;
      for (let k of Object.keys(lockedBalances)) {
        let tokenInfoState = state.getIn(['tokensInfo', accountAddress, k]);
        if (!tokenInfoState) {
          tokenInfoState = initialTokenInfo;
        }
        tokenInfoState = tokenInfoState.set('lockedBalance', lockedBalances[k]);
        state = state.setIn(['tokensInfo', accountAddress, k], tokenInfoState);
      }
      return state;
    }
    case 'LOGIN':
      state = state.setIn(['isLoggedIn', action.payload.address], true);
      return state;
    case 'LOGOUT':
      state = state.setIn(['isLoggedIn', action.payload.address], false);
      return state;
    case 'LOAD_ORDERS':
      state = state.set('orders', OrderedMap());
      action.payload.orders.reverse().forEach(o => {
        state = state.setIn(['orders', o.id], o);
      });
      return state;
    case 'ORDER_UPDATE':
      const order = action.payload.order;
      const ordersPath = ['orders', order.id];

      if (state.getIn(ordersPath)) {
        if (order.status !== 'pending') {
          state = state.deleteIn(ordersPath);
        } else {
          state = state.setIn(ordersPath, order);
        }
      } else if (order.status === 'pending') {
        state = state.setIn(ordersPath, order);
      }
      return state;
    case 'CANCEL_ORDER':
      state = state.deleteIn(['orders', action.payload.id]);
      return state;
    case 'LOAD_TRADES':
      state = state.set('trades', OrderedMap());
      action.payload.trades.reverse().forEach(t => {
        state = state.setIn(['trades', t.id], t);
      });
      return state;
    case 'TRADE_UPDATE':
      let trade = action.payload.trade;
      state = state.setIn(['trades', trade.id], trade);
      return state;
    case 'UPDATE_TOKEN_INFO': {
      const { symbol, balance, allowance, decimals, tokenAddress, accountAddress } = action.payload;
      let tokenInfoState = state.getIn(['tokensInfo', accountAddress, symbol]);
      if (!tokenInfoState) {
        tokenInfoState = initialTokenInfo;
      }
      tokenInfoState = tokenInfoState.set('allowance', allowance);
      tokenInfoState = tokenInfoState.set('balance', balance);
      tokenInfoState = tokenInfoState.set('address', tokenAddress);
      if (decimals) {
        tokenInfoState = tokenInfoState.set('decimals', decimals);
      }
      state = state.setIn(['tokensInfo', accountAddress, symbol], tokenInfoState);
      return state;
    }
    default:
      return state;
  }
};
