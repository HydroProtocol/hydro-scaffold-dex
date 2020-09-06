import React from 'react';
import { connect } from 'react-redux';
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

  maxAmount(asks) {
    let max = 0;
    asks.forEach(element => {
      const amount = element[1].toFixed(this.props.currentMarket.amountDecimals);
      if (max < amount) {
        max = amount;
      }
    });
    return max;
  }

  calculateBarWidth(maxAmount, Amount) {
    const width = ~~((Amount / maxAmount) * 20); //20 is the max percentage
    return width + '%';
  }

  render() {
    const { bids, asks, websocketConnected, currentMarket, dollarExchangeRate } = this.props;
    const asksArray = asks.slice(-20).reverse().toArray();
    const bidsArray = bids.slice(-20).reverse().toArray();
    const asksMaxAmount = this.maxAmount(asksArray);
    const bidsMaxAmount = this.maxAmount(bidsArray);

    return (
      <div className="orderbook">
        <div className="flex header text-secondary">
          <div className="col-4 text-center border-right font-weight-bold">Amount</div>
          <div className="col-4 text-center border-right font-weight-bold">Price</div>
          <div className="col-4 text-center font-weight-bold">Amount</div>
        </div>
        <div className="flex-column flex-1">
          <div className="asks flex-column flex-column-reverse">
            {asksArray
              .map(([price, amount], index) => {
                const dollarValue = price * dollarExchangeRate;
                return (
                  <div className="ask flex align-items-center" style={{backgroundColor: index % 2 ? '#f7f7f7' : '#fff' }} key={price.toString()}>
                    <div className="col-4 orderbook-amount text-center">
                      <p className="amount">{amount.toFixed(currentMarket.amountDecimals)}</p>
                      <div className="amount-bar-ask" style={{width: this.calculateBarWidth(asksMaxAmount, amount.toFixed(currentMarket.amountDecimals))}}></div>
                    </div>
                    <div className="col-4 text-center font-weight-bold price-container-ask">
                      <div className="price text-danger">{price.toFixed(currentMarket.priceDecimals)}</div>
                      <div className="currency">{dollarValue.toFixed(2)} USD</div>
                    </div>
                    <div className="col-4 text-center">-</div>
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
          <div className="bids flex-column">
            {bidsArray
              .map(([price, amount], index) => {
                const dollarValue = price * dollarExchangeRate;
                return (
                  <div className="ask flex align-items-center" style={{backgroundColor: index % 2 ? '#f7f7f7' : '#fff' }} key={price.toString()}>
                    <div className="col-4 text-center">-</div>
                    <div className="col-4 text-center font-weight-bold bid-price-container">
                      <div className="price text-primary">{price.toFixed(currentMarket.priceDecimals)}</div>
                      <div className="currency">{dollarValue.toFixed(2)} USD</div>
                    </div>
                    <div className="col-4 orderbook-amount text-center">
                      <div className="amount-bar-bid" style={{width: this.calculateBarWidth(bidsMaxAmount, amount.toFixed(currentMarket.amountDecimals))}}></div>
                      <p className="amount">{amount.toFixed(currentMarket.amountDecimals)}</p>
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
  const currentMarket = state.market.getIn(['markets', 'currentMarket']);
  return {
    asks: state.market.getIn(['orderbook', 'asks']),
    bids: state.market.getIn(['orderbook', 'bids']),
    loading: false,
    currentMarket,
    websocketConnected: state.config.get('websocketConnected'),
    theme: state.config.get('theme'),
    dollarExchangeRate: state.market.getIn(['exchangeRate', 'data', currentMarket['quoteToken']]),
  };
};

export default connect(mapStateToProps)(OrderBook);
