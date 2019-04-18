import BigNumber from 'bignumber.js';
import { hotDiscountRules } from '../actions/fee';

// a pure function to caculate all trade details
export const calculateTrade = ({
  orderType,
  side,
  price,
  amount,
  hotTokenAmount,
  gasFeeAmount,
  asMakerFeeRate,
  asTakerFeeRate,
  amountDecimals,
  priceDecimals
}) => {
  let tradeFee, subtotal, totalBaseTokens;

  const isMakerFee = orderType === 'limit';
  const hotDiscount = getHotDiscountRate(hotTokenAmount);
  const feeRate = orderType === 'market' ? asTakerFeeRate : asMakerFeeRate;

  if (orderType === 'market' && side === 'buy') {
    subtotal = amount;
    totalBaseTokens = amount.div(price).dp(amountDecimals, BigNumber.ROUND_DOWN);
  } else {
    subtotal = amount.multipliedBy(price).dp(priceDecimals, BigNumber.ROUND_DOWN);
    totalBaseTokens = amount;
  }

  const estimatedPrice = orderType === 'market' ? price : new BigNumber(0);

  tradeFee = subtotal.multipliedBy(feeRate);
  const tradeFeeAfterDiscount = tradeFee.multipliedBy(hotDiscount);
  const feeRateAfterDiscount = feeRate.multipliedBy(hotDiscount);
  let totalQuoteTokens;
  if (side === 'buy') {
    totalQuoteTokens = subtotal
      .plus(tradeFeeAfterDiscount)
      .plus(gasFeeAmount)
      .dp(priceDecimals, BigNumber.ROUND_UP);
  } else {
    totalQuoteTokens = subtotal
      .minus(tradeFeeAfterDiscount)
      .minus(gasFeeAmount)
      .dp(priceDecimals, BigNumber.ROUND_UP);
  }
  totalQuoteTokens = BigNumber.max(totalQuoteTokens, new BigNumber('0'));

  return {
    estimatedPrice,
    gasFeeAmount,
    hotDiscount,
    totalBaseTokens,
    tradeFeeAfterDiscount,
    feeRateAfterDiscount,
    tradeFee,
    feeRate,
    isMakerFee,
    subtotal,
    totalQuoteTokens
  };
};

const getHotDiscountRate = hotTokenAmount => {
  hotTokenAmount = hotTokenAmount.div(10 ** 18);
  for (let rule of hotDiscountRules) {
    const limit = new BigNumber(rule[0]);
    const discountRate = new BigNumber(rule[1]);

    if (limit.eq(-1)) {
      return discountRate;
    } else if (hotTokenAmount.lte(limit)) {
      return discountRate;
    }
  }

  return new BigNumber(1);
};
