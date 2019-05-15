import BigNumber from 'bignumber.js';
import { watchToken } from '../actions/account';
import abi from './abi';
import env from './env';
import { getSelectedAccountWallet, getTransactionReceipt, getContract } from '@gongddex/hydro-sdk-wallet';
export let web3, Contract;

export const getTokenBalance = async (tokenAddress, accountAddress) => {
  const contract = getContract(tokenAddress, abi);
  const balance = await contract.balanceOf(accountAddress);
  return new BigNumber(balance);
};

export const getAllowance = async (tokenAddress, accountAddress) => {
  const contract = getContract(tokenAddress, abi);
  const allowance = await contract.allowance(accountAddress, env.HYDRO_PROXY_ADDRESS);
  return new BigNumber(allowance);
};

export const wrapETH = amount => {
  return async (dispatch, getState) => {
    const state = getState();
    const WETH = state.config.get('WETH');
    const value = new BigNumber(amount).multipliedBy(Math.pow(10, WETH.decimals)).toString();

    let params = {
      to: WETH.address,
      data: '0xd0e30db0',
      value
    };

    try {
      const wallet = getSelectedAccountWallet(state);
      const transactionID = await wallet.sendTransaction(params);

      alert(`Wrap ETH request submitted`);
      watchTransactionStatus(transactionID, async success => {
        if (success) {
          dispatch(watchToken(WETH.address, WETH.symbol, WETH.decimals));
          alert('Wrap ETH Successfully');
        } else {
          alert('Wrap ETH Failed');
        }
      });
      return transactionID;
    } catch (e) {
      alert(e);
    }
    return null;
  };
};

export const unwrapWETH = amount => {
  return async (dispatch, getState) => {
    const state = getState();
    const WETH = state.config.get('WETH');
    const value = new BigNumber(amount).multipliedBy(Math.pow(10, WETH.decimals)).toString(16);
    const wallet = getSelectedAccountWallet(state);
    const functionSelector = '2e1a7d4d';
    const valueString = get64BytesString(value);

    let params = {
      to: WETH.address,
      data: `0x${functionSelector}${valueString}`,
      value: 0
    };

    try {
      const transactionID = await wallet.sendTransaction(params);

      alert(`Unwrap WETH request submitted`);
      watchTransactionStatus(transactionID, async success => {
        if (success) {
          dispatch(watchToken(WETH.address, WETH.symbol, WETH.decimals));
          alert('Wrap ETH Successfully');
        } else {
          alert('Wrap ETH Failed');
        }
      });
      return transactionID;
    } catch (e) {
      alert(e);
    }
    return null;
  };
};

export const enable = (address, symbol, decimals) => {
  return async (dispatch, getState) => {
    let transactionID = await dispatch(
      approve(address, symbol, 'f000000000000000000000000000000000000000000000000000000000000000', 'Enable', decimals)
    );
    return transactionID;
  };
};

export const disable = (address, symbol, decimals) => {
  return async (dispatch, getState) => {
    let transactionID = await dispatch(
      approve(address, symbol, '0000000000000000000000000000000000000000000000000000000000000000', 'Disable', decimals)
    );
    return transactionID;
  };
};

export const approve = (tokenAddress, symbol, allowance, action, decimals) => {
  return async (dispatch, getState) => {
    const state = getState();
    const functionSelector = '095ea7b3';
    let spender = get64BytesString(env.HYDRO_PROXY_ADDRESS);
    if (spender.length !== 64) {
      return null;
    }

    let params = {
      to: tokenAddress,
      data: `0x${functionSelector}${spender}${allowance}`,
      value: 0
    };

    try {
      const wallet = getSelectedAccountWallet(state);
      const transactionID = await wallet.sendTransaction(params);

      alert(`${action} ${symbol} request submitted`);
      watchTransactionStatus(transactionID, async success => {
        if (success) {
          dispatch(watchToken(tokenAddress, symbol, decimals));
          alert(`${action} ${symbol} Successfully`);
        } else {
          alert(`${action} ${symbol} Failed`);
        }
      });
      return transactionID;
    } catch (e) {
      alert(e);
    }
    return null;
  };
};

const watchTransactionStatus = (txID, callback) => {
  const getTransactionStatus = async () => {
    const tx = await getTransactionReceipt(txID);
    if (!tx) {
      window.setTimeout(() => getTransactionStatus(txID), 3000);
    } else if (callback) {
      callback(Number(tx.status) === 1);
    } else {
      alert('success');
    }
  };
  window.setTimeout(() => getTransactionStatus(txID), 3000);
};

const get64BytesString = string => {
  string = string.replace('0x', '');
  while (string.length < 64) {
    string = '0'.concat(string);
  }
  return string;
};
