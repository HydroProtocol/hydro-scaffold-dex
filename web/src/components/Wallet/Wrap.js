import React from 'react';
import { connect } from 'react-redux';
import BigNumber from 'bignumber.js';
import { wrapETH, unwrapWETH } from '../../lib/wallet';
import { toUnitAmount } from '../../lib/utils';
import { stateUtils } from '../../selectors/account';
import { getSelectedAccount } from '@gongddex/hydro-sdk-wallet';

const mapStateToProps = state => {
  const WETH = state.config.get('WETH');
  const selectedAccount = getSelectedAccount(state);
  const ethBalance = selectedAccount ? selectedAccount.get('balance') : new BigNumber('0');
  const address = selectedAccount ? selectedAccount.get('address') : null;
  const wethBalance = stateUtils.getTokenAvailableBalance(state, address, 'WETH');
  return {
    ethBalance: toUnitAmount(ethBalance, 18),
    wethBalance: toUnitAmount(wethBalance, WETH.decimals)
  };
};

class Wrap extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      amount: ''
    };
  }

  componentDidUpdate(prevProps) {
    const { type } = this.props;
    if (type !== prevProps.type) {
      this.setState({ amount: '' });
    }
  }

  render() {
    const { ethBalance, wethBalance, type } = this.props;
    const { amount } = this.state;
    const isWrap = type === 'wrap';

    return (
      <form className="form flex-column text-secondary flex-1 justify-content-between block">
        <div className="form-group">
          <label className="text-secondary">
            Amount ({isWrap ? ethBalance.toFixed(8) : wethBalance.toFixed(8)} Max)
          </label>
          <div className="input-group">
            <input
              className="form-control"
              value={amount}
              onChange={event => this.setState({ amount: event.target.value })}
            />
          </div>
        </div>
        <button
          type="button"
          className={`form-control btn ${isWrap ? 'btn-success' : 'btn-danger'}`}
          onClick={() => this.submit()}>
          {type}
        </button>
      </form>
    );
  }

  submit() {
    const { dispatch, type } = this.props;
    const { amount } = this.state;
    if (type === 'wrap') {
      dispatch(wrapETH(amount));
    } else {
      dispatch(unwrapWETH(amount));
    }
  }
}

export default connect(mapStateToProps)(Wrap);
