package api

import (
	"encoding/json"
	"fmt"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type BuildOrderTestSuit struct {
	suite.Suite
}

func (suite *BuildOrderTestSuit) SetupTest() {
	setEnvs()
	models.InitTestDBPG()

	mockMarketDao()
	mockErc20()
	mockLockedBlanceDao()
	mockCacheService()
}

func mockLockedBlanceDao() {
	balanceDao := models.MLockedBalanceDao{}
	models.BalanceDao = &balanceDao
	balanceDao.On("GetByAccountAndSymbol", mock.Anything, mock.Anything, mock.Anything).Times(10).Return(decimal.Zero)
}

func mockMarketDao() {
	marketDao := &models.MMarketDao{}

	models.MarketDao = marketDao
	var markets []*models.Market

	marketWethDai := &models.Market{
		ID:                 "WETH-DAI",
		BaseTokenSymbol:    "WETH",
		BaseTokenAddress:   os.Getenv("HSK_WETH_TOKEN_ADDRESS"),
		BaseTokenDecimals:  18,
		QuoteTokenSymbol:   "DAI",
		QuoteTokenAddress:  os.Getenv("HSK_USD_TOKEN_ADDRESS"),
		QuoteTokenDecimals: 18,
		MinOrderSize:       decimal.NewFromFloat(0.1),
		PricePrecision:     5,
		PriceDecimals:      5,
		AmountDecimals:     5,
		MakerFeeRate:       decimal.NewFromFloat(0.001),
		TakerFeeRate:       decimal.NewFromFloat(0.003),
		GasUsedEstimation:  250000,
	}

	marketHotDai := &models.Market{
		ID:                 "HOT-DAI",
		BaseTokenSymbol:    "HOT",
		BaseTokenAddress:   os.Getenv("HSK_WETH_TOKEN_ADDRESS"),
		BaseTokenDecimals:  18,
		QuoteTokenSymbol:   "DAI",
		QuoteTokenAddress:  os.Getenv("HSK_USD_TOKEN_ADDRESS"),
		QuoteTokenDecimals: 18,
		MinOrderSize:       decimal.NewFromFloat(0.1),
		PricePrecision:     5,
		PriceDecimals:      5,
		AmountDecimals:     5,
		MakerFeeRate:       decimal.NewFromFloat(0.001),
		TakerFeeRate:       decimal.NewFromFloat(0.003),
		GasUsedEstimation:  250000,
	}

	markets = append(markets, marketWethDai)
	markets = append(markets, marketHotDai)

	marketDao.On("FindAllMarkets").Return(markets).Once()
	marketDao.On("FindMarketByID", mock.MatchedBy(func(marketID string) bool { return marketID == "WETH-DAI" })).Return(marketWethDai)
	marketDao.On("FindMarketByID", "HOT-DAI").Times(10).Return(marketHotDai)
	marketDao.On("FindMarketByID", mock.AnythingOfType("string")).Times(10).Return(nil)
}

func mockErc20() {
	mockHydro := sdk.MockHydro{
		&ethereum.EthereumHydroProtocol{},
		&sdk.MockBlockchain{},
	}

	hydro = mockHydro

	aBigNumber, _ := decimal.NewFromString("10000000000000000000000000000000")

	mockHydro.BlockChain.(*sdk.MockBlockchain).On("GetTokenBalance", mock.Anything, mock.Anything).Times(10).Return(aBigNumber, nil)
	mockHydro.BlockChain.(*sdk.MockBlockchain).On("GetTokenAllowance", mock.Anything, mock.Anything, mock.Anything).Times(10).Return(aBigNumber, nil)
	mockHydro.BlockChain.(*sdk.MockBlockchain).On("GetHotFeeDiscount", mock.Anything).Times(10).Return(decimal.New(1, 2))
	mockHydro.BlockChain.(*sdk.MockBlockchain).On("IsValidSignature", mock.Anything, mock.Anything, mock.Anything).Return(true, nil)
}

func mockCacheService() {
	cacheService := &models.MCache{}

	cacheService.On("Set", mock.Anything, mock.Anything, mock.Anything).Times(100).Return(nil)
	cacheService.On("Get", mock.Anything).Times(10).Return("", nil)
	cacheService.On("Push", mock.Anything).Times(10).Return(nil)
	cacheService.On("Pop").Times(1000).Return([]string{}, nil)

	CacheService = cacheService
	QueueService = cacheService
}

func (suite *BuildOrderTestSuit) request(body interface{}) *Response {
	e := getEchoServer()
	bts, _ := json.Marshal(body)
	fmt.Println("request is ", string(bts))
	req := httptest.NewRequest(http.MethodPost, "/orders/build", strings.NewReader(string(bts)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	address := "0x5409ed021d9299bf6814279a6a1411a7e866a631"
	signatrue := "0xdcd19ecc53c51bc1c8c67183d9ed8a2c68bb3717b7bbbd39da969960feeb95d45f79ead1d476c5cb1f2ebf77b76a87abee2bf5643a235125a85428d3ef4926b700"
	message := "HYDRO-AUTHENTICATION"
	req.Header.Set("Hydro-Authentication", address+"#"+message+"#"+signatrue)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	var res Response
	json.Unmarshal(rec.Body.Bytes(), &res)
	return &res
}

func (suite *BuildOrderTestSuit) TestInvalidParams() {
	res := suite.request(map[string]interface{}{})
	suite.Equal(-1, res.Status)
}

func (suite *BuildOrderTestSuit) TestValidParams() {
	req := BuildOrderReq{
		BaseReq:   BaseReq{Address: "0x5409ed021d9299bf6814279a6a1411a7e866a631"},
		MarketID:  "HOT-DAI",
		Side:      "buy",
		OrderType: "limit",
		Price:     "140",
		Amount:    "100",
	}

	res := suite.request(req)
	fmt.Println(res.Desc)
	fmt.Println(utils.ToJsonString(res.Data))
	suite.Equal(0, res.Status)
}

func TestBuildOrder(t *testing.T) {
	suite.Run(t, new(BuildOrderTestSuit))
}
