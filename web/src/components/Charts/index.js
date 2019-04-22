import React from 'react';
import { connect } from 'react-redux';
import { DeepChart, TradeChart } from '@wangleiddex/hydro-sdk-charts';
import { testData } from './constants';

class Charts extends React.Component {
  render() {
    const bids = this.props.bids.toArray().map(priceLevel => {
      return {
        price: priceLevel[0].toString(),
        amount: priceLevel[1].toString()
      };
    });
    const asks = this.props.asks.toArray().map(priceLevel => {
      return {
        price: priceLevel[0].toString(),
        amount: priceLevel[1].toString()
      };
    });

    return (
      <>
        <div className="title flex justify-content-between align-items-center">
          <div>
            <div>Charts</div>
          </div>
        </div>

        <div className="flex-column flex-1 ">
          <div className="grid flex-2">
            <TradeChart
              theme="light"
              data={testData}
              priceDecimals={5}
              styles={{ upColor: '#00d99f', downColor: '#ff6f75', background: '#FFFFFF' }}
              clickCallback={result => {
                console.log('result: ', result);
              }}
              handleLoadMore={result => {
                console.log('result: ', result);
              }}
              clickGranularity={result => {
                console.log('result: ', result);
              }}
            />
          </div>
          <div className="grid flex-1 border-top">
            <DeepChart
              baseToken="HOT"
              quoteToken="DAI"
              styles={{ bidColor: '#00d99f', askColor: '#ff6f75', rowBackgroundColor: '#FFFFFF' }}
              asks={asks}
              bids={bids}
              priceDecimals={5}
              theme="light"
              clickCallback={result => {
                console.log('result: ', result);
              }}
            />
          </div>
        </div>
      </>
    );
  }
}

const mapStateToProps = state => {
  return {
    asks: state.market.getIn(['orderbook', 'asks']),
    bids: state.market.getIn(['orderbook', 'bids'])
  };
};

export default connect(mapStateToProps)(Charts);
