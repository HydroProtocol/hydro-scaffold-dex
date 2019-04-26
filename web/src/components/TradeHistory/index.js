import React from 'react';
import { connect } from 'react-redux';
import BigNumber from 'bignumber.js';
import PerfectScrollbar from 'perfect-scrollbar';
import moment from 'moment';

const mapStateToProps = state => {
  return {
    tradeHistory: state.market.get('tradeHistory'),
    currentMarket: state.market.getIn(['markets', 'currentMarket'])
  };
};

class TradeHistory extends React.PureComponent {
  componentDidUpdate(prevProps) {
    const { tradeHistory } = this.props;
    if (tradeHistory !== prevProps.tradeHistory) {
      this.ps.update();
    }
  }

  render() {
    const { tradeHistory, currentMarket } = this.props;
    return (
      <div className="trade-history flex-1 position-relative overflow-hidden" ref={ref => this.setRef(ref)}>
        <table className="table">
          <thead>
            <tr className="text-secondary">
              <th className="text-right">Price</th>
              <th className="text-right">Amount</th>
              <th>Time</th>
            </tr>
          </thead>
          <tbody>
            {tradeHistory
              .toArray()
              .reverse()
              .map(([id, trade]) => {
                const colorGreen = trade.takerSide === 'buy';
                return (
                  <tr key={trade.id}>
                    <td className={['text-right', colorGreen ? 'text-success' : 'text-danger'].join(' ')}>
                      {new BigNumber(trade.price).toFixed(currentMarket.priceDecimals)}
                      {trade.takerSide === 'buy' ? (
                        <i className="fa fa-arrow-up" aria-hidden="true" />
                      ) : (
                        <i className="fa fa-arrow-down" aria-hidden="true" />
                      )}
                    </td>
                    <td className="text-right">{new BigNumber(trade.amount).toFixed(currentMarket.amountDecimals)}</td>
                    <td className="text-secondary">{moment(trade.executedAt).format('hh:mm:ss')}</td>
                  </tr>
                );
              })}
          </tbody>
        </table>
      </div>
    );
  }

  setRef(ref) {
    if (ref) {
      this.ps = new PerfectScrollbar(ref, {
        suppressScrollX: true,
        maxScrollbarLength: 20
      });
    }
  }
}

export default connect(mapStateToProps)(TradeHistory);
