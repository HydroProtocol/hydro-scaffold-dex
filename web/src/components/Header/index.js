import React from 'react';
import { loginRequest, login } from '../../actions/account';
import { updateCurrentMarket } from '../../actions/markets';
import { connect } from 'react-redux';
import { Wallet, WalletButton, getSelectedAccount } from '@gongddex/hydro-sdk-wallet';
import './styles.scss';
import { loadAccountHydroAuthentication } from '../../lib/session';
import env from '../../lib/env';

const mapStateToProps = state => {
  const selectedAccountID = state.WalletReducer.get('selectedAccountID');
  const selectedAccount = getSelectedAccount(state);
  const address = selectedAccount ? selectedAccount.get('address') : null;
  return {
    address,
    isLocked: selectedAccount ? selectedAccount.get('isLocked') : true,
    isLoggedIn: state.account.getIn(['isLoggedIn', address]),
    currentMarket: state.market.getIn(['markets', 'currentMarket']),
    markets: state.market.getIn(['markets', 'data']),
    networkId: state.WalletReducer.getIn(['accounts', selectedAccountID, 'networkId'])
  };
};

class Header extends React.PureComponent {
  componentDidMount() {
    const { address, dispatch } = this.props;
    const hydroAuthentication = loadAccountHydroAuthentication(address);
    if (hydroAuthentication) {
      dispatch(login(address));
    }
  }
  componentDidUpdate(prevProps) {
    const { address, dispatch } = this.props;
    const hydroAuthentication = loadAccountHydroAuthentication(address);
    if (address !== prevProps.address && hydroAuthentication) {
      dispatch(login(address));
    }
  }
  render() {
    const { currentMarket, markets, dispatch, networkId } = this.props;
    return (
      <div className="navbar bg-blue navbar-expand-lg">
        <img className="navbar-brand" src={require('../../images/hydro.svg')} alt="hydro" />
        <div className="dropdown navbar-nav mr-auto">
          <button
            className="btn btn-primary header-dropdown dropdown-toggle"
            type="button"
            id="dropdownMenuButton"
            data-toggle="dropdown"
            aria-haspopup="true"
            aria-expanded="false">
            {currentMarket && currentMarket.id}
          </button>
          <div
            className="dropdown-menu"
            aria-labelledby="dropdownMenuButton"
            style={{ maxHeight: 350, overflow: 'auto' }}>
            {markets.map(market => {
              return (
                <button
                  className="dropdown-item"
                  key={market.id}
                  onClick={() => currentMarket.id !== market.id && dispatch(updateCurrentMarket(market))}>
                  {market.id}
                </button>
              );
            })}
          </div>
        </div>
        {parseInt(networkId, 10) !== 66 && (
          <span className="btn text-danger" style={{ marginRight: 12 }}>
            Network Error: Switch Metamask's network to localhost:8545.
          </span>
        )}
        <a
          href="https://hydroprotocol.io/developers/docs/overview/what-is-hydro.html"
          className="btn btn-primary"
          target="_blank"
          rel="noopener noreferrer"
          style={{ marginRight: 12 }}>
          DOCUMENTATION
        </a>
        <WalletButton />
        <Wallet title="Starter Kit Wallet" nodeUrl={env.NODE_URL} defaultWalletType="Browser Wallet" />
        {this.renderAccount()}
      </div>
    );
  }

  renderAccount() {
    const { address, dispatch, isLoggedIn, isLocked } = this.props;
    if ((isLoggedIn && address) || isLocked) {
      return null;
    } else if (address) {
      return (
        <button className="btn btn-success" style={{ marginLeft: 12 }} onClick={() => dispatch(loginRequest())}>
          connect
        </button>
      );
    } else {
      return null;
    }
  }
}

export default connect(mapStateToProps)(Header);
