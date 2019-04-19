package admincli

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/admin/api"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"math/big"
	"os"
	"testing"
	"time"
)

type AdminTest struct {
	suite.Suite
	admin IAdminApi
}

func (a *AdminTest) SetupSuite() {
}

func (a *AdminTest) SetupTest() {
	utils.Error("sfsdfsf%s", "f")
	os.Setenv("ADMIN_API_URL", "")
	os.Setenv("HSK_REDIS_URL", "redis://127.0.0.1:6379/0")
	os.Setenv("HSK_DATABASE_URL", "")
	os.Setenv("HSK_HYBRID_EXCHANGE_ADDRESS", "")
	os.Setenv("WEB_HEALTH_CHECK_URL", "")
	os.Setenv("API_HEALTH_CHECK_URL", "")
	os.Setenv("ENGINE_HEALTH_CHECK_URL", "")
	os.Setenv("LAUNCHER_HEALTH_CHECK_URL", "")
	os.Setenv("WATCHER_HEALTH_CHECK_URL", "")
	os.Setenv("WEBSOCKET_HEALTH_CHECK_URL", "")
	os.Setenv("HSK_BLOCKCHAIN_RPC_URL", "http://localhost:8545")
	os.Setenv("HSK_HYBRID_EXCHANGE_ADDRESS", "")

	go adminapi.StartServer(context.Background())
	<-time.After(time.Second)
	a.admin = NewAdmin("", nil, &MockErc20{})
}

func (a *AdminTest) TestStatus() {
	assert.Nil(a.T(), a.admin.Status())
}

func (a *AdminTest) TestNewMarket() {
	var marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation string
	assert.Nil(a.T(), a.admin.NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation))
}

func (a *AdminTest) TestUpdateMarket() {
	var marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish string
	assert.Nil(a.T(), a.admin.UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish))
}

func (a *AdminTest) TestPublishMarket() {
	var marketID string
	assert.Nil(a.T(), a.admin.PublishMarket(marketID))
}

func (a *AdminTest) TestUnPublishMarket() {
	var marketID string
	assert.Nil(a.T(), a.admin.UnPublishMarket(marketID))
}

func (a *AdminTest) TestUpdateMarketFee() {
	var marketID, makerFee, takerFee string
	assert.Nil(a.T(), a.admin.UpdateMarketFee(marketID, makerFee, takerFee))
}

func (a *AdminTest) TestListAccountOrders() {
	var address, limit, offset, status string
	assert.Nil(a.T(), a.admin.ListAccountOrders(address, limit, offset, status))
}

func (a *AdminTest) TestListAccountBalances() {
	var address, limit, offset string
	assert.Nil(a.T(), a.admin.ListAccountBalances(address, limit, offset))
}

func (a *AdminTest) TestListAccountTrades() {
	var address, limit, offset, status string
	assert.Nil(a.T(), a.admin.ListAccountTrades(address, limit, offset, status))
}

func TestAdmin(t *testing.T) {
	suite.Run(t, new(AdminTest))
}

type MockErc20 struct {
}

func (MockErc20) Symbol(address string) (error, string) {
	return nil, "HOT"
}

func (MockErc20) Decimals(address string) (error, int) {
	return nil, 18
}

func (MockErc20) Name(address string) (error, string) {
	return nil, "hydro token"
}

func (MockErc20) TotalSupply(address string) (error, *big.Int) {
	return nil, big.NewInt(10000)
}
