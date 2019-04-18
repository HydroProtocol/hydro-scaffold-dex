import React from 'react';
import { connect } from 'react-redux';
import { loadTokens } from '../../actions/account';
import { toUnitAmount, isTokenApproved } from '../../lib/utils';
import BigNumber from 'bignumber.js';
import { enable, disable } from '../../lib/wallet';
import { getSelectedAccount } from 'hydro-sdk-wallet';

const mapStateToProps = state => {
  const selectedType = state.WalletReducer.get('selectedType');
  const selectedAccount = getSelectedAccount(state);
  const address = selectedAccount ? selectedAccount.get('address') : null;
  return {
    tokensInfo: state.account.get('tokensInfo'),
    address,
    lockedBalances: state.account.get('lockedBalances'),
    isLoggedIn: state.account.getIn(['isLoggedIn', address]),
    ethBalance: toUnitAmount(state.WalletReducer.getIn(['accounts', selectedType, 'balance']), 18)
  };
};

class Tokens extends React.PureComponent {
  componentDidMount() {
    const { address, dispatch, isLoggedIn } = this.props;
    if (address && isLoggedIn) {
      dispatch(loadTokens());
    }
  }

  componentDidUpdate(prevProps) {
    const { address, dispatch, isLoggedIn } = this.props;
    const accountChange = address !== prevProps.address;
    const loggedInChange = isLoggedIn !== prevProps.isLoggedIn;
    if (address && isLoggedIn && (accountChange || loggedInChange)) {
      dispatch(loadTokens());
    }
  }

  render() {
    const { dispatch, tokensInfo, lockedBalances, ethBalance } = this.props;
    return (
      <div className="flex-column">
        <div className="token flex flex-1">
          <div className="col-6">ETH</div>
          <div className="col-6 text-right">{ethBalance.toFixed(5)}</div>
        </div>
        {tokensInfo.toArray().map(([token, info]) => {
          const { address, balance, allowance, decimals } = info.toJS();
          const lockedBalance = lockedBalances.get(token, new BigNumber('0'));
          const isApproved = isTokenApproved(allowance || new BigNumber('0'));
          const availableBalance = toUnitAmount(balance.minus(lockedBalance) || new BigNumber('0'), decimals).toFixed(
            5
          );
          const toolTipTitle = `<div>In-Order: ${toUnitAmount(lockedBalance, decimals).toFixed(
            5
          )}</div><div>Total: ${toUnitAmount(balance, decimals).toFixed(5)}</div>`;
          return (
            <div key={token} className="token flex flex-1">
              <div className="flex-column col-6">
                <div>{token}</div>
                <div className="text-secondary">{isApproved ? 'Enabled' : 'Disabled'}</div>
              </div>
              <div className="col-6 text-right">
                <div
                  className="flex-column"
                  key={toolTipTitle}
                  data-html="true"
                  data-toggle="tooltip"
                  data-placement="right"
                  title={toolTipTitle}
                  ref={ref => window.$(ref).tooltip()}>
                  {availableBalance}
                </div>
                <div className="custom-control custom-switch">
                  <input
                    type="checkbox"
                    className="custom-control-input"
                    id={address}
                    checked={isApproved}
                    onChange={() => {
                      if (isApproved) {
                        dispatch(disable(address, token));
                      } else {
                        dispatch(enable(address, token));
                      }
                    }}
                  />
                  <label className="custom-control-label" htmlFor={address} />
                </div>
              </div>
            </div>
          );
        })}
      </div>
    );
  }
}

export default connect(mapStateToProps)(Tokens);
