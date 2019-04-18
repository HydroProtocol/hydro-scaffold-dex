export const initOrderbook = (bids, asks) => {
  return async dispatch => {
    dispatch({
      type: 'INIT_ORDERBOOK',
      payload: {
        bids,
        asks
      }
    });
  };
};

export const updateOrderbook = (side, price, amount) => {
  return dispatch => {
    return dispatch({
      type: 'UPDATE_ORDERBOOK',
      payload: {
        side,
        price,
        amount
      }
    });
  };
};
