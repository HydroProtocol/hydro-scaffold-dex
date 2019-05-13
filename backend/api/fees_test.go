package api

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestFees(t *testing.T) {
	models.MockMarketDao()

	mockBlockChain := &sdk.MockBlockchain{}
	mockBlockChain.On("GetHotFeeDiscount", mock.Anything).Return(decimal.New(1, 0))
	mockHydro := sdk.MockHydro{
		&ethereum.EthereumHydroProtocol{},
		mockBlockChain,
	}

	hydro = mockHydro

	market := models.MarketDao.FindMarketByID("HOT-DAI")
	// decimal18
	decimals18 := decimal.New(1, 18)
	// 3gwei
	gasPrice := decimal.New(3, 9)
	// ethInDai
	ethPriceInDai := decimal.NewFromFloat(150)
	// gas
	gasFeeAmt := decimal.NewFromFloat(float64(market.GasUsedEstimation)).Mul(gasPrice).Mul(ethPriceInDai).Div(decimals18)
	asMakerTotalFee := decimal.NewFromFloat(140 * 100).Mul(market.MakerFeeRate).Add(gasFeeAmt)
	asTakerTotalFee := decimal.NewFromFloat(140 * 100).Mul(market.TakerFeeRate).Add(gasFeeAmt)

	url := "/fees?price=140&amount=100&marketID=HOT-DAI"
	resp := request(url, "GET", "", nil)

	assert.EqualValues(t, 0, resp.Status)
	fees := resp.Data.(map[string]interface{})["fees"]
	//spew.Dump(resp.Data)
	assert.EqualValues(t, "0.001", fees.(map[string]interface{})["asMakerFeeRate"])
	assert.EqualValues(t, "14", fees.(map[string]interface{})["asMakerTradeFeeAmount"])
	assert.EqualValues(t, asMakerTotalFee.String(), fees.(map[string]interface{})["asMakerTotalFeeAmount"])
	assert.EqualValues(t, "0.003", fees.(map[string]interface{})["asTakerFeeRate"])
	assert.EqualValues(t, "42", fees.(map[string]interface{})["asTakerTradeFeeAmount"])
	assert.EqualValues(t, asTakerTotalFee.String(), fees.(map[string]interface{})["asTakerTotalFeeAmount"])
	assert.EqualValues(t, gasFeeAmt.String(), fees.(map[string]interface{})["gasFeeAmount"])
}

// api should return -4 when amount or price is not positive number
func TestInvalidRequest(t *testing.T) {
	models.MockMarketDao()
	url := "/fees?price=140&amount=0&marketID=HOT-DAI"
	resp := request(url, "GET", "", nil)

	assert.EqualValues(t, -4, resp.Status)
}

// not exist market AAA-CCC
func TestUnExitMarket(t *testing.T) {
	models.MockMarketDao()
	url := "/fees?price=140&amount=0&marketID=AAA-CCC"
	resp := request(url, "GET", "", nil)

	assert.EqualValues(t, -3, resp.Status)
}
