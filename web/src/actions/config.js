export const setConfigs = configs => {
  return {
    type: 'SET_CONFIGS',
    payload: configs
  };
};

export const loadWeb3NetworkID = web3NetworkID => {
  web3NetworkID = parseInt(web3NetworkID, 10);
  return async (dispatch, getState) => {
    const state = getState();
    const oldWeb3NetworkID = state.config.get('web3NetworkID');

    if (oldWeb3NetworkID === web3NetworkID) {
      return;
    } else if (oldWeb3NetworkID) {
      window.location.reload();
    }

    dispatch(setConfigs({ web3NetworkID }));
  };
};
