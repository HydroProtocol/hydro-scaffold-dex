import React from 'react';
import { connect } from 'react-redux';
import { loadTokens } from '../../actions/account';
import { toUnitAmount, isTokenApproved } from '../../lib/utils';
import { stateUtils } from '../../selectors/account';
import { enable, disable } from '../../lib/wallet';
import { getSelectedAccount } from '@gongddex/hydro-sdk-wallet';
import { BigNumber } from 'bignumber.js';

const mapStateToProps = state => {
  const selectedAccountID = state.WalletReducer.get('selectedAccountID');
  const selectedAccount = getSelectedAccount(state);
  const address = selectedAccount ? selectedAccount.get('address') : null;
  return {
    tokensInfo: stateUtils.getTokensInfo(state, address),
    address,
    ethBalance: toUnitAmount(state.WalletReducer.getIn(['accounts', selectedAccountID, 'balance']), 18)
  };
};

class Tokens extends React.PureComponent {
  componentDidMount() {
    const { address, dispatch } = this.props;
    if (address) {
      dispatch(loadTokens());
    }
  }

  componentDidUpdate(prevProps) {
    const { address, dispatch } = this.props;
    if (address && address !== prevProps.address) {
      dispatch(loadTokens());
    }
  }

  render() {
    const { dispatch, tokensInfo, ethBalance } = this.props;
    return (
      <div className="flex-column">
        <div className="token flex flex-1">
          <div className="col-6">ETH</div>
          <div className="col-6 text-right">{ethBalance.toFixed(5)}</div>
        </div>
        {tokensInfo.toArray().map(([token, info]) => {
          const { address, balance, allowance, decimals, lockedBalance } = info.toJS();
          const isApproved = isTokenApproved(allowance);
          const availableBalance = toUnitAmount(BigNumber.max(balance.minus(lockedBalance), '0'), decimals).toFixed(5);
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
                        dispatch(disable(address, token, decimals));
                      } else {
                        dispatch(enable(address, token, decimals));
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
