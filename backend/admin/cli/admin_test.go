package admincli

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/admin/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type AdminTest struct {
	admin IAdminApi
}

func (a AdminTest) T() *testing.T {
	panic("implement me")
}

func (a AdminTest) SetT(*testing.T) {
	panic("implement me")
}

func (a AdminTest) SetupTest() {
	os.Setenv("ADMIN_API_URL", "")
	os.Setenv("HSK_DATABASE_URL", "")
	os.Setenv("HSK_HYBRID_EXCHANGE_ADDRESS", "")
	os.Setenv("WEB_HEALTH_CHECK_URL", "")
	os.Setenv("API_HEALTH_CHECK_URL", "")
	os.Setenv("ENGINE_HEALTH_CHECK_URL", "")
	os.Setenv("LAUNCHER_HEALTH_CHECK_URL", "")
	os.Setenv("WATCHER_HEALTH_CHECK_URL", "")
	os.Setenv("WEBSOCKET_HEALTH_CHECK_URL", "")
	os.Setenv("HSK_BLOCKCHAIN_RPC_URL", "")
	os.Setenv("HSK_HYBRID_EXCHANGE_ADDRESS", "")

	adminapi.StartServer(context.Background())
	a.admin = NewAdmin("")
}

func (a *AdminTest) TestStatus(t *testing.T) {

	a.admin.Status()
}

func (a *AdminTest) TestNewMarket(t *testing.T) {
	var marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation string
	assert.Nil(t, admin.NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation))
}

func (a *AdminTest) TestUpdateMarket(t *testing.T) {

}

func (a *AdminTest) TestPublishMarket(t *testing.T) {

}

func (a *AdminTest) TestUnPublishMarket(t *testing.T) {

}

func (a *AdminTest) TestUpdateMarketFee(t *testing.T) {

}

func (a *AdminTest) TestListAccountOrders(t *testing.T) {

}

func (a *AdminTest) TestListAccountBalances(t *testing.T) {

}

func (a *AdminTest) TestListAccountTrades(t *testing.T) {
}

func TestAdmin(t *testing.T) {
	suite.Run(t, new(AdminTest))
}
