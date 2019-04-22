package admincli

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/admin/api"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/davecgh/go-spew/spew"
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
	_, err := a.admin.Status()
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestNewMarket() {
	var marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation string
	marketID = "HOT-WETH"
	baseTokenAddress = "0x0000000000000000000000000000000000000001"
	quoteTokenAddress = "0x0000000000000000000000000000000000000002"
	_, err := a.admin.NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestUpdateMarket() {
	var marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish string
	_, err := a.admin.UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestPublishMarket() {
	var marketID string
	_, err := a.admin.PublishMarket(marketID)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestUnPublishMarket() {
	var marketID string
	_, err := a.admin.UnPublishMarket(marketID)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestUpdateMarketFee() {
	var marketID, makerFee, takerFee string
	_, err := a.admin.UpdateMarketFee(marketID, makerFee, takerFee)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestListAccountOrders() {
	var address, limit, offset, status string
	_, err := a.admin.ListAccountOrders(address, limit, offset, status)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestListAccountBalances() {
	var address, limit, offset string
	_, err := a.admin.ListAccountBalances(address, limit, offset)
	assert.Nil(a.T(), err)
}

func (a *AdminTest) TestListAccountTrades() {
	var address, limit, offset, status string
	_, err := a.admin.ListAccountTrades(address, limit, offset, status)
	assert.Nil(a.T(), err)
}

func TestAdmin(t *testing.T) {
	suite.Run(t, new(AdminTest))
}

//mock erc20 service
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

func TestArray(t *testing.T) {

	a := []string{"0", "1", "2", "3", "4", "5"}
	spew.Dump(a[0:1])
	spew.Dump(a[0:0])
	spew.Dump(a[1:1])
	spew.Dump(a[0:2])
	spew.Dump(a[2:5])
	spew.Dump(a[2:6])
}
