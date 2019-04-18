import { Map, OrderedMap } from 'immutable';
import BigNumber from 'bignumber.js';

export const initState = Map({
  address: null,
  isLoggedIn: Map(),
  ethBalance: new BigNumber('0'),
  lockedBalances: Map(),
  tokensInfo: Map(),
  approving: Map(),
  orders: OrderedMap(),
  trades: OrderedMap(),
  transactions: OrderedMap()
});

export default (state = initState, action) => {
  switch (action.type) {
    case 'LOAD_ACCOUNT':
      state = state.set('address', action.payload.address);
      return state;
    case 'LOAD_BALANCE':
      state = state.set('ethBalance', action.payload.balance);
      return state;
    case 'UPDATE_TOKEN_LOCKED_BALANCES':
      for (let k of Object.keys(action.payload)) {
        state = state.setIn(['lockedBalances', k], action.payload[k]);
      }
      return state;
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
    case 'LOAD_TOKEN':
      const { symbol, balance, allowance, decimals, address } = action.payload;
      state = state.setIn(['tokensInfo', symbol, 'allowance'], allowance);
      state = state.setIn(['tokensInfo', symbol, 'balance'], balance);
      state = state.setIn(['tokensInfo', symbol, 'address'], address);
      if (decimals) {
        state = state.setIn(['tokensInfo', symbol, 'decimals'], decimals);
      }
      return state;
    default:
      return state;
  }
};
