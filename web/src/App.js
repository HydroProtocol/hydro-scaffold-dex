import React from 'react';
import { connect } from 'react-redux';
import { loadMarkets, loadTradeHistory } from './actions/markets';
import Header from './components/Header';
import WebsocketConnector from './components/WebsocketConnector';
import OrderBook from './components/Orderbook';
import Trade from './components/Trade';
import Wallet from './components/Wallet';
import Orders from './components/Orders';
import TradeHistory from './components/TradeHistory';

const mapStateToProps = state => {
  return {
    currentMarket: state.market.getIn(['markets', 'currentMarket'])
  };
};

class App extends React.PureComponent {
  componentDidMount() {
    const { dispatch, currentMarket } = this.props;
    dispatch(loadMarkets());
    if (currentMarket) {
      dispatch(loadTradeHistory(currentMarket.id));
    }
  }

  componentDidUpdate(prevProps) {
    const { currentMarket, dispatch } = this.props;
    if (currentMarket !== prevProps.currentMarket) {
      dispatch(loadTradeHistory(currentMarket.id));
    }
  }

  render() {
    const { currentMarket } = this.props;
    if (!currentMarket) {
      return null;
    }
    return (
      <div className="app">
        <WebsocketConnector />
        <Header />
        <div className="flex flex-1 overflow-hidden">
          <div className="flex">
            <div className="grid border-right">
              <Trade />
            </div>
            <div className="grid border-right flex-column">
              <OrderBook />
            </div>
          </div>
          <div className="flex-column flex-1 border-right">
            <div className="grid flex-1">
              <Wallet />
            </div>
            <div className="grid flex-1 border-top">
              <Orders />
            </div>
          </div>
          <div className="grid">
            <TradeHistory />
          </div>
        </div>
      </div>
    );
  }
}

export default connect(mapStateToProps)(App);
