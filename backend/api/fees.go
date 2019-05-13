package api

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
)

func GetFees(p Param) (interface{}, error) {
	params := p.(*FeesReq)

	price := utils.StringToDecimal(params.Price)
	amount := utils.StringToDecimal(params.Amount)

	market := models.MarketDao.FindMarketByID(params.MarketID)

	if market == nil {
		return nil, MarketNotFoundError(params.MarketID)
	}

	quoteTokenTotalAmount := price.Mul(amount)

	if quoteTokenTotalAmount.Equal(decimal.Zero) {
		return nil, InvalidPriceAmountError()
	}

	fee := calculateFee(price, amount, market, params.Address)

	fees := FeesResp{
		GasFeeAmount: fee.GasFeeAmount,

		AsMakerTradeFeeAmount: fee.AsMakerTradeFeeAmount,
		AsMakerTotalFeeAmount: fee.AsMakerTotalFeeAmount,
		AsMakerFeeRate:        market.MakerFeeRate,

		AsTakerTradeFeeAmount: fee.AsTakerTradeFeeAmount,
		AsTakerTotalFeeAmount: fee.AsTakerTotalFeeAmount,
		AsTakerFeeRate:        market.TakerFeeRate,
	}

	return map[string]interface{}{
		"fees": fees,
	}, nil
}

type feeDetail struct {
	Address     string
	HotDiscount decimal.Decimal

	Price          decimal.Decimal
	Amount         decimal.Decimal
	AsMakerFeeRate decimal.Decimal
	AsTakerFeeRate decimal.Decimal

	GasFeeAmount decimal.Decimal

	AsMakerTradeFeeAmount decimal.Decimal
	AsMakerTotalFeeAmount decimal.Decimal

	AsTakerTradeFeeAmount decimal.Decimal
	AsTakerTotalFeeAmount decimal.Decimal
}

// Hydro Relayer is in charge of sending transaction to Ethereum network.
// The traders have to pay some gas fee to cover Relayer gas cost.
// Otherwise, a trader can place a batch of orders with very small amount to make the Relayer run out of ether.
// To provide minimum permission request for traders, the Relayer charges gas fee in quote token.
// For a normal ERC20 token, a match on hydro will cost about 180000 - 190000 GWei gas on Ethereum.
//
// For an example, a trader is placing an order for HOT-DAI market.
// Let's assume:
//   1) ETH price is 150 DAI
//   2) Reasonable gas price is 10 GWei
// So the gas fee should be 150 * 180000 * 0.000000001 * 10 = 0.27 DAI
func getGasFeeAmount(market *models.Market) decimal.Decimal {
	// 3GWei
	gasPrice := decimal.New(3, 9)

	gasCostInEth := gasPrice.Mul(decimal.NewFromFloat(float64(market.GasUsedEstimation)))

	var gasCostInQuoteToken decimal.Decimal
	if market.QuoteTokenSymbol == "WETH" {
		gasCostInQuoteToken = gasCostInEth
	} else {
		// for markets with other quote token, assume its DAI for simplicity, should replace this logic in real world app

		// assume WETH's price is: 150 DAI
		wethPriceInDai := decimal.NewFromFloat(150)

		gasCostInQuoteToken = gasCostInEth.Mul(wethPriceInDai)
	}

	humanFriendlyGasCostInQuoteToken := gasCostInQuoteToken.Div(decimal.New(1, 18))
	return humanFriendlyGasCostInQuoteToken
}

func calculateFee(price, amount decimal.Decimal, market *models.Market, address string) *feeDetail {
	detail := &feeDetail{}

	detail.Price = price
	detail.Amount = amount
	detail.HotDiscount = hydro.GetHotFeeDiscount(address)

	detail.AsMakerFeeRate = market.MakerFeeRate
	detail.AsTakerFeeRate = market.TakerFeeRate
	detail.GasFeeAmount = getGasFeeAmount(market)

	detail.AsMakerTradeFeeAmount = market.MakerFeeRate.Mul(price).Mul(amount).Mul(detail.HotDiscount)
	detail.AsMakerTotalFeeAmount = detail.AsMakerTradeFeeAmount.Add(detail.GasFeeAmount)

	detail.AsTakerTradeFeeAmount = market.TakerFeeRate.Mul(price).Mul(amount).Mul(detail.HotDiscount)
	detail.AsTakerTotalFeeAmount = detail.AsTakerTradeFeeAmount.Add(detail.GasFeeAmount)

	return detail
}
