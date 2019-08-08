import React from 'react';
import { connect } from 'react-redux';
import BarChart from '../BarChart'
import './styles.scss';

class OrderBook extends React.Component {
  constructor(props) {
    super(props);
    this.lastUpdatedAt = null;
    this.forceRenderTimer = null;
  }

  // max 1 render in 1 second
  shouldComponentUpdate() {
    if (this.lastUpdatedAt) {
      const diff = new Date().valueOf() - this.lastUpdatedAt;
      const shouldRender = diff > 1000;

      if (!shouldRender && !this.forceRenderTimer) {
        this.forceRenderTimer = setTimeout(() => {
          this.forceUpdate();
          this.forceRenderTimer = null;
        }, 1000 - diff);
      }
      return shouldRender;
    } else {
      return true;
    }
  }

  componentWillUnmount() {
    if (this.forceRenderTimer) {
      clearInterval(this.forceRenderTimer);
    }
  }

  componentDidUpdate() {
    this.lastUpdatedAt = new Date();
  }

  render() {
    let { bids, asks, websocketConnected, currentMarket } = this.props;
    console.log('his.props',this.props);

    return (
      <div className="orderbook flex-column flex-1">
        <div className="flex header text-secondary">
          <div className="col-4 text-center"> <b>Amount</b> </div>
          <div className="col-4 text-center"><b>Price</b></div>
          <div className="col-4 text-center">-  </div>
        </div>
        <div className="flex-column flex-1">
          <div className="asks flex-column flex-column-reverse overflow-hidden">
            {asks
              .slice(-20)
              .reverse()
              .toArray()
              .map(([price, amount], index) => {
                const barSize = (price/amount ) * 10000; // Some confusion how to calculate percentage  
                return (
                  <div className={`ask flex align-items-center ${index%2 && 'orderbook--zebraGray'}`} key={price.toString()}>
                    <div className="col-4 orderbook-amount text-left">
                      <BarChart percent={barSize}/>
                      <div className="orderbook-amount-value">{amount.toFixed(currentMarket.amountDecimals)} </div>
                    </div>
                    <div className="col-4 text-danger text-center orderbook--opacityGray"><div><b>{price.toFixed(currentMarket.priceDecimals)}</b>
                    </div><div className="orderbook--currency">79 <b>USD</b></div>
                    </div>
                    <div className="col-4 orderbook-amount text-center">
                    -
                    </div>
                  </div>
                );
              })}
          </div>
          <div className="status border-top border-bottom">
            {websocketConnected ? (
              <div className="col-6 text-success">
                <i className="fa fa-circle" aria-hidden="true" /> RealTime
              </div>
            ) : (
              <div className="col-6 text-danger">
                <i className="fa fa-circle" aria-hidden="true" /> Disconnected
              </div>
            )}
          </div>
          <div className="bids flex-column flex-1 overflow-hidden">
            {bids
              .slice(0, 20)
              .toArray()
              .map(([price, amount], index) => {
                const barSize = (price/amount ) * 10000; 
                return (
                  <div className={`bid flex align-items-center ${index%2 && 'orderbook--zebraGray'}`} key={price.toString()}>
                    <div className="col-4 orderbook-amount text-center">
                    <BarChart percent={barSize}/>
                      <div className="orderbook-amount-value">{amount.toFixed(currentMarket.amountDecimals)}</div>
                    </div>
                    <div className="col-4 text-success text-center orderbook--opacityGray"><div><b>{price.toFixed(currentMarket.priceDecimals)}</b></div>
                    <div className="orderbook--currency">79 <b>USD</b></div>
                    </div>
                    <div className="col-4 text-center">
                    -
                    </div>
                  </div>
                );
              })}
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    asks: state.market.getIn(['orderbook', 'asks']),
    bids: state.market.getIn(['orderbook', 'bids']),
    loading: false,
    currentMarket: state.market.getIn(['markets', 'currentMarket']),
    websocketConnected: state.config.get('websocketConnected'),
    theme: state.config.get('theme')
  };
};

export default connect(mapStateToProps)(OrderBook);
