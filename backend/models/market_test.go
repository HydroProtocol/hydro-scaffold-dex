package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_PG_GetDefaultMarketDao(t *testing.T) {
	setEnvs()
	InitTestDBPG()

	markets := MarketDaoPG.FindAllMarkets()
	assert.EqualValues(t, 0, len(markets))
}

func Test_PG_MarketDao_FindAndInsertMarket(t *testing.T) {
	setEnvs()
	InitTestDBPG()

	dbMarket := MarketDaoPG.FindMarketByID("HOT-WETH")
	assert.Nil(t, dbMarket)

	market := Market{
		ID:                 "HOT-WETH",
		BaseTokenSymbol:    "HOT",
		BaseTokenName:      "HOT",
		QuoteTokenSymbol:   "WETH",
		QuoteTokenName:     "WETH",
		BaseTokenAddress:   os.Getenv("HSK_HYDRO_TOKEN_ADDRESS"),
		QuoteTokenAddress:  os.Getenv("HSK_WETH_TOKEN_ADDRESS"),
		BaseTokenDecimals:  18,
		QuoteTokenDecimals: 18,
		MinOrderSize:       utils.IntToDecimal(1),
		PricePrecision:     8,
		PriceDecimals:      8,
		AmountDecimals:     8,
		MakerFeeRate:       utils.StringToDecimal("0.001"),
		TakerFeeRate:       utils.StringToDecimal("0.001"),
		GasUsedEstimation:  250000,
		IsPublished:        true,
	}

	assert.Nil(t, MarketDaoPG.InsertMarket(&market))
	dbMarket = MarketDaoPG.FindMarketByID("HOT-WETH")

	publishMarkets := MarketDaoPG.FindPublishedMarkets()
	assert.EqualValues(t, 1, len(publishMarkets))

	assert.EqualValues(t, market.ID, dbMarket.ID)
	assert.EqualValues(t, market.BaseTokenDecimals, dbMarket.BaseTokenDecimals)
	assert.EqualValues(t, market.BaseTokenSymbol, dbMarket.BaseTokenSymbol)
	assert.EqualValues(t, market.BaseTokenName, dbMarket.BaseTokenName)
}
