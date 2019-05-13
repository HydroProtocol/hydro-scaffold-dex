import React from 'react';
import { connect } from 'react-redux';
import { formValueSelector, Field, stopSubmit } from 'redux-form';
import { TRADE_FORM_ID } from '../../actions/trade';
import { reduxForm } from 'redux-form';
import { trade } from '../../actions/trade';
import BigNumber from 'bignumber.js';
import { loadHotDiscountRules, getHotTokenAmount } from '../../actions/fee';
import { calculateTrade } from '../../lib/tradeCalculator';
import { loginRequest } from '../../actions/account';
import PerfectScrollbar from 'perfect-scrollbar';
import './styles.scss';
import { sleep, toUnitAmount } from '../../lib/utils';
import { getSelectedAccount } from '@gongddex/hydro-sdk-wallet';
import { stateUtils } from '../../selectors/account';

const mapStateToProps = state => {
  const selector = formValueSelector(TRADE_FORM_ID);
  const bids = state.market.getIn(['orderbook', 'bids']);
  const asks = state.market.getIn(['orderbook', 'asks']);
  const selectedAccount = getSelectedAccount(state);
  const address = selectedAccount ? selectedAccount.get('address') : null;
  const currentMarket = state.market.getIn(['markets', 'currentMarket']);
  const lastTrade = state.market.get('tradeHistory').first();
  const lastPrice = lastTrade ? new BigNumber(lastTrade.price) : new BigNumber('0');

  return {
    initialValues: {
      side: 'buy',
      orderType: 'limit',
      subtotal: new BigNumber(0),
      total: new BigNumber(0),
      totalBase: new BigNumber(0),
      feeRate: new BigNumber(0),
      gasFee: new BigNumber(0),
      hotDiscount: new BigNumber(1),
      tradeFee: new BigNumber(0),
      estimatedPrice: new BigNumber(0),
      marketOrderWorstPrice: new BigNumber(0),
      marketOrderWorstTotalQuote: new BigNumber(0),
      marketOrderWorstTotalBase: new BigNumber(0)
    },
    lastPrice,
    currentMarket,
    quoteTokenBalance: stateUtils.getTokenAvailableBalance(state, address, currentMarket.quoteToken),
    baseTokenBalance: stateUtils.getTokenAvailableBalance(state, address, currentMarket.baseToken),
    hotTokenAmount: state.config.get('hotTokenAmount'),
    address,
    isLoggedIn: state.account.getIn(['isLoggedIn', address]),
    price: new BigNumber(selector(state, 'price') || 0),
    amount: new BigNumber(selector(state, 'amount') || 0),
    total: new BigNumber(selector(state, 'total') || 0),
    totalBase: new BigNumber(selector(state, 'totalBase') || 0),
    subtotal: new BigNumber(selector(state, 'subtotal') || 0),
    feeRate: new BigNumber(selector(state, 'feeRate') || 0),
    gasFee: new BigNumber(selector(state, 'gasFee') || 0),
    estimatedPrice: new BigNumber(selector(state, 'estimatedPrice') || 0),
    marketOrderWorstPrice: new BigNumber(selector(state, 'marketOrderWorstPrice') || 0),
    marketOrderWorstTotalQuote: new BigNumber(selector(state, 'marketOrderWorstTotalQuote') || 0),
    marketOrderWorstTotalBase: new BigNumber(selector(state, 'marketOrderWorstTotalBase') || 0),
    hotDiscount: new BigNumber(selector(state, 'hotDiscount') || 1),
    tradeFee: new BigNumber(selector(state, 'tradeFee') || 0),
    side: selector(state, 'side'),
    orderType: selector(state, 'orderType'),
    bestBidPrice: bids.size > 0 ? bids.get(0)[0].toString() : null,
    bestAskPrice: asks.size > 0 ? asks.get(asks.size - 1)[0].toString() : null
  };
};

class Trade extends React.PureComponent {
  componentDidMount() {
    const { dispatch } = this.props;
    loadHotDiscountRules();
    this.interval = window.setInterval(() => {
      dispatch(getHotTokenAmount());
    }, 30 * 1000);
  }

  componentDidUpdate(prevProps) {
    const { currentMarket, reset, lastPrice, price, change } = this.props;
    if (currentMarket.id === prevProps.currentMarket.id) {
      if (!lastPrice.eq(prevProps.lastPrice) && price.eq(0)) {
        change('price', lastPrice);
      }
      this.updateFees(prevProps);
    } else {
      reset();
    }
  }

  render() {
    const { side, handleSubmit, currentMarket, total, gasFee, tradeFee, subtotal, change } = this.props;
    if (!currentMarket) {
      return null;
    }

    return (
      <>
        <div className="title">
          <div>
            <div>{currentMarket.id}</div>
            <div className="text-secondary">Make a Limit Order</div>
          </div>
        </div>
        <div className="trade flex-1 flex-column">
          <ul className="nav nav-tabs">
            <li className="nav-item flex-1 flex">
              <div
                className={`flex-1 tab-button text-secondary text-center${side === 'buy' ? ' active' : ''}`}
                onClick={() => change('side', 'buy')}>
                Buy
              </div>
            </li>
            <li className="nav-item flex-1 flex">
              <div
                className={`flex-1 tab-button text-secondary text-center${side === 'sell' ? ' active' : ''}`}
                onClick={() => change('side', 'sell')}>
                Sell
              </div>
            </li>
          </ul>
          <div className="flex flex-1 position-relative overflow-hidden" ref={ref => this.setRef(ref)}>
            <form
              className="form flex-column text-secondary flex-1 justify-content-between"
              onSubmit={handleSubmit(() => this.submit())}>
              <div>
                <Field
                  name="price"
                  unit={currentMarket.quoteToken}
                  autoComplete="off"
                  component={this.renderField}
                  label="Price"
                />
                <Field
                  name="amount"
                  unit={currentMarket.baseToken}
                  autoComplete="off"
                  component={this.renderField}
                  label="Amount"
                />
                <div className="form-group">
                  <div className="form-title">Order Summary</div>
                  <div className="list">
                    <div className="item flex justify-content-between">
                      <div className="name">Order</div>
                      <div className="name">{subtotal.toFixed(currentMarket.priceDecimals)}</div>
                    </div>
                    <div className="item flex justify-content-between">
                      <div className="name">Fees</div>
                      <div className="name">{gasFee.plus(tradeFee).toFixed(currentMarket.priceDecimals)}</div>
                    </div>
                    <div className="item flex justify-content-between">
                      <div className="name">Total</div>
                      <div className="name">{total.toFixed(currentMarket.priceDecimals)}</div>
                    </div>
                  </div>
                </div>
              </div>
              <button type="submit" className={`form-control btn ${side === 'buy' ? 'btn-success' : 'btn-danger'}`}>
                {side} {currentMarket.baseToken}
              </button>
            </form>
          </div>
        </div>
      </>
    );
  }

  renderField = ({ input, label, unit, meta, ...attrs }) => {
    const { submitFailed, error } = meta;

    return (
      <div className="form-group">
        <label>{label}</label>
        <div className="input-group">
          <input className="form-control" {...input} {...attrs} />
          <span className="text-secondary unit">{unit}</span>
        </div>
        <span className="text-danger">{submitFailed && (error && <span>{error}</span>)}</span>
      </div>
    );
  };

  async submit() {
    const { amount, price, side, orderType, dispatch, isLoggedIn, address } = this.props;
    if (!isLoggedIn) {
      await dispatch(loginRequest(address));
      // Metamask's window will be hidden when continuous call Metamask sign method
      await sleep(500);
    }
    try {
      await dispatch(trade(side, price, amount, orderType));
    } catch (e) {
      alert(e);
    }
  }

  updateFees(prevProps) {
    const { currentMarket, orderType, side, price, amount, hotTokenAmount, change } = this.props;

    if (
      orderType === prevProps.orderType &&
      side === prevProps.side &&
      price.eq(prevProps.price) &&
      amount.eq(prevProps.amount) &&
      hotTokenAmount.eq(prevProps.hotTokenAmount)
    ) {
      return;
    }
    const { asMakerFeeRate, asTakerFeeRate, gasFeeAmount, priceDecimals, amountDecimals } = currentMarket;

    const calculateParam = {
      orderType,
      side,
      price: new BigNumber(price),
      amount: new BigNumber(amount),
      hotTokenAmount,
      gasFeeAmount,
      asMakerFeeRate,
      asTakerFeeRate,
      amountDecimals,
      priceDecimals
    };

    const calculateResult = calculateTrade(calculateParam);

    change('subtotal', calculateResult.subtotal);
    change('estimatedPrice', calculateResult.estimatedPrice);
    change('totalBase', calculateResult.totalBaseTokens);
    change('total', calculateResult.totalQuoteTokens);
    change('feeRate', calculateResult.feeRateAfterDiscount);
    change('gasFee', calculateResult.gasFeeAmount);
    change('hotDiscount', calculateResult.hotDiscount);
    change('tradeFee', calculateResult.tradeFeeAfterDiscount);
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

const validate = (values, props) => {
  const { price, amount, total } = values;
  const { side, address, currentMarket, quoteTokenBalance, baseTokenBalance } = props;

  let _price, _amount, _total;

  const errors = {};

  if (address) {
    if (side === 'buy') {
      const quoteTokenAmount = toUnitAmount(quoteTokenBalance, currentMarket.quoteTokenDecimals);

      if (quoteTokenAmount.eq(0)) {
        errors.amount = `Insufficient ${currentMarket.quoteToken} balance`;
      }
    } else {
      const baseTokenAmount = toUnitAmount(baseTokenBalance, currentMarket.baseTokenDecimals);
      if (baseTokenAmount.eq(0)) {
        errors.amount = `Insufficient ${currentMarket.baseToken} balance`;
      }
    }
  }

  if (!errors.price) {
    if (!price) {
      errors.price = 'Price required';
    } else if (isNaN(Number(price))) {
      errors.price = 'Price must be a number';
    } else {
      _price = new BigNumber(price);
      if (_price.lte('0')) {
        errors.price = `Price cannot be 0`;
      }
    }
  }
  if (!errors.amount) {
    if (!amount) {
      errors.amount = 'Amount required';
    } else if (isNaN(Number(amount))) {
      errors.amount = 'Amount must be a number';
    } else {
      _amount = new BigNumber(amount);

      if (_amount.lte('0')) {
        errors.amount = `Amount cannot be 0`;
      } else if (_amount.multipliedBy(_price).lt(currentMarket.minOrderSize)) {
        errors.amount = `total sale price too small`;
      }
    }
  }

  if (!errors.amount && !errors.price && total && address) {
    _total = new BigNumber(total);
    if (side === 'buy') {
      const quoteTokenAmount = toUnitAmount(quoteTokenBalance, currentMarket.quoteTokenDecimals);

      if (_total.gt(quoteTokenAmount)) {
        errors.amount = `Insufficient ${currentMarket.quoteToken} balance`;
      }
    } else {
      const baseTokenAmount = toUnitAmount(baseTokenBalance, currentMarket.baseTokenDecimals);

      if (_amount.gt(baseTokenAmount)) {
        errors.amount = `Insufficient ${currentMarket.baseToken} balance`;
      } else if (_total.lte('0')) {
        errors.amount = `Amount too small: total sale price less than fee`;
      }
    }
  }
  return errors;
};

const shouldError = () => {
  return true;
};
const onSubmitFail = (_, dispatch) => {
  setTimeout(() => {
    dispatch(stopSubmit(TRADE_FORM_ID));
  }, 3000);
};

export default connect(mapStateToProps)(
  reduxForm({
    form: TRADE_FORM_ID,
    destroyOnUnmount: false,
    onSubmitFail,
    validate,
    shouldError
  })(Trade)
);
