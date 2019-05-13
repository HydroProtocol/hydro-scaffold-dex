package dex_engine

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/engine"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"math/rand"
	"os"
	"testing"
	"time"
)

type marketHandlerSuite struct {
	suite.Suite
	marketHandler *MarketHandler
}

const fakeAccount1 = "0x31ebd457b999bf99759602f5ece5aa5033cb56b3"
const fakeAccount2 = "0x3eb06f432ae8f518a957852aa44776c234b4a84a"

func (s *marketHandlerSuite) SetupSuite() {
}

func (s *marketHandlerSuite) TearDownSuite() {
}

func (s *marketHandlerSuite) SetupTest() {
	setEnvs()
	models.InitTestDBPG()
	//models.MockMarketDao()
	market := &models.Market{
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

	err := models.MarketDao.InsertMarket(market)
	if err != nil {
		panic(err)
	}

	token := &models.Token{
		Name:     "HOT",
		Symbol:   "HOT",
		Address:  os.Getenv("HSK_WETH_TOKEN_ADDRESS"),
		Decimals: 18,
	}

	_ = models.TokenDao.InsertToken(token)

	wsQueue = &common.MockQueue{}
	kvStore := &common.MockKVStore{}

	wsQueue.(*common.MockQueue).On("Push", mock.Anything).Return(nil)
	kvStore.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	kvStore.On("Get", mock.Anything).Return("", common.KVStoreEmpty)
	marketHotDai := models.MarketHotDai()
	marketHandler, _ := NewMarketHandler(context.Background(), marketHotDai, engine.NewEngine(context.Background()))
	s.marketHandler = marketHandler

	s.marketHandler.hydroEngine.RegisterOrderBookActivitiesHandler(RedisOrderBookActivitiesHandler{})
}

func (s *marketHandlerSuite) TearDownTest() {
}

func (s *marketHandlerSuite) AssertChange(realTestFunc func(), counterFunc func() int, expectedDiff int, msgAndArgs ...interface{}) {
	resBefore := counterFunc()
	realTestFunc()
	resAfter := counterFunc()
	s.Equal(expectedDiff, resAfter-resBefore, msgAndArgs...)
}

func (s *marketHandlerSuite) TestHandleNewOrder() {
	sellOrder := newModelOrder("buy", decimal.New(140, 0), decimal.New(140, 0))

	sellOrderEvent := &common.NewOrderEvent{
		Event: common.Event{},
		Order: utils.ToJsonString(sellOrder),
	}

	//s.Nil(s.marketHandler.orderbook.MaxBid())

	s.AssertChange(func() {
		_, _ = s.marketHandler.handleNewOrder(sellOrderEvent)
	}, func() int {
		return models.OrderDao.Count()
	}, 1)

	//s.Equal("140", s.marketHandler.orderbook.MaxBid().String())
}

type batchMatchOrdersTest struct {
	takerOrderParams          *buildOrderParams
	makerOrdersParams         []*buildOrderParams
	expectedTradesCount       int
	expectedTransactionsCount int

	whenPending *expectedResult
	whenSuccess *expectedResult
	whenFailed  *expectedResult

	takerOrder  *models.Order
	makerOrders []*models.Order
}

func (b *batchMatchOrdersTest) Reset() {
	b.takerOrder = b.takerOrderParams.toModelOrder()
	b.makerOrders = make([]*models.Order, 0, 10)

	for i := range b.makerOrdersParams {
		params := b.makerOrdersParams[i]
		b.makerOrders = append(b.makerOrders, params.toModelOrder())
	}

}

type buildOrderParams struct {
	side   string
	price  string
	amount string
}

func (p *buildOrderParams) toModelOrder() *models.Order {
	return newModelOrder(
		p.side,
		utils.StringToDecimal(p.price),
		utils.StringToDecimal(p.amount),
	)
}

type expectedResult struct {
	expectedAmounts               [][]string
	expectedStatus                []string
	expectedMarketChannelPayloads []*common.WebsocketMarketOrderChangePayload
}

func contains(s [][]byte, e []byte) bool {
	for _, a := range s {
		if string(a) == string(e) {
			return true
		}
	}
	return false
}

func (s *marketHandlerSuite) newBatchMatchOrdersTest(
	takerOrderParams *buildOrderParams,
	makerOrdersParams []*buildOrderParams,
	expectedTradesCount int,
	expectedTransactionsCount int,
	whenPending *expectedResult,
	whenSuccess *expectedResult,
	whenFailed *expectedResult,
) {
	testConfig := &batchMatchOrdersTest{
		takerOrderParams,
		makerOrdersParams,
		expectedTradesCount,
		expectedTransactionsCount,
		whenPending,
		whenSuccess,
		whenFailed,
		nil,
		nil,
	}

	s.batchNewOrderTest(testConfig)
}

func (s *marketHandlerSuite) batchNewOrderTest(b *batchMatchOrdersTest) {
	b.Reset()
	s.batchNewOrderTestPendingPart(b)

	if b.whenSuccess != nil {
		s.SetupTest()
		b.Reset()
		_, launchLog := s.batchNewOrderTestPendingPart(b)
		hash := "fake-success"
		launchLog.Hash = sql.NullString{
			hash,
			true,
		}
		models.UpdateLaunchLogToPending(launchLog)
		takerOrderEvent := common.ConfirmTransactionEvent{
			Event:  common.Event{},
			Hash:   hash,
			Status: common.STATUS_SUCCESSFUL,
		}
		_, _ = s.marketHandler.handleTransactionResult(&takerOrderEvent)
		s.assertExpectedResult(b, b.whenSuccess)
	}

	if b.whenFailed != nil {
		s.SetupTest()
		b.Reset()
		_, launchLog := s.batchNewOrderTestPendingPart(b)
		hash := "fake-failed"
		launchLog.Hash = sql.NullString{
			hash,
			true,
		}
		models.UpdateLaunchLogToPending(launchLog)
		takerOrderEvent := common.ConfirmTransactionEvent{
			Event:  common.Event{},
			Hash:   hash,
			Status: common.STATUS_FAILED,
		}
		_, _ = s.marketHandler.handleTransactionResult(&takerOrderEvent)
		s.assertExpectedResult(b, b.whenFailed)
	}
}

func (s *marketHandlerSuite) assertExpectedResult(b *batchMatchOrdersTest, result *expectedResult) {
	// reload orders
	b.takerOrder = models.OrderDao.FindByID(b.takerOrder.ID)
	for i := range b.makerOrders {
		b.makerOrders[i] = models.OrderDao.FindByID(b.makerOrders[i].ID)
	}

	for i, status := range result.expectedStatus {
		if i == 0 {
			s.Equal(status, b.takerOrder.Status)
		} else {
			s.Equal(status, b.makerOrders[i-1].Status)
		}
	}

	for i, amounts := range result.expectedAmounts {
		if i == 0 {
			s.assertOrderAmounts(amounts[0], amounts[1], amounts[2], amounts[3], b.takerOrder)
		} else {
			s.assertOrderAmounts(amounts[0], amounts[1], amounts[2], amounts[3], b.makerOrders[i-1])
		}
	}

	queueBuffers := wsQueue.(*common.MockQueue).Buffers

	if result.expectedMarketChannelPayloads != nil {
		for i := range result.expectedMarketChannelPayloads {
			payload := result.expectedMarketChannelPayloads[i]

			msg := &common.WebSocketMessage{
				ChannelID: common.GetMarketChannelID(s.marketHandler.market.ID),
				Payload:   payload,
			}

			//expectedMsg, _ := json.Marshal(msg)
			//log.Println("msg expected:", string(expectedMsg))
			//for _, real := range queueBuffers {
			//	log.Println(" == ", string(real))
			//}

			msgBytes, _ := json.Marshal(msg)
			s.True(contains(queueBuffers, msgBytes), fmt.Sprintf("msg %s not exist", msgBytes))
		}
	}

	assertHasOrderChangeMsgFunc := func(order *models.Order) {
		msg := common.WebSocketMessage{
			ChannelID: common.GetAccountChannelID(order.TraderAddress),
			Payload: &common.WebsocketOrderChangePayload{
				Type:  common.WsTypeOrderChange,
				Order: order,
			},
		}
		msgBytes, _ := json.Marshal(msg)
		s.True(contains(queueBuffers, msgBytes), fmt.Sprintf("msg %s not exist", msgBytes))
	}

	// There must be some order change events
	assertHasOrderChangeMsgFunc(b.takerOrder)
	for i := range b.makerOrders {
		makerOrder := b.makerOrders[i]
		assertHasOrderChangeMsgFunc(makerOrder)
	}
}

func (s *marketHandlerSuite) batchNewOrderTestPendingPart(b *batchMatchOrdersTest) (*models.Transaction, *models.LaunchLog) {
	oldTradesCount := models.TradeDao.Count()
	oldTransactionsCount := models.TransactionDao.Count()

	for _, makerOrder := range b.makerOrders {
		makerOrderEvent := common.NewOrderEvent{
			Event: common.Event{},
			Order: utils.ToJsonString(makerOrder),
		}

		_, _ = s.marketHandler.handleNewOrder(&makerOrderEvent)
	}

	takerOrderEvent := common.NewOrderEvent{
		Event: common.Event{},
		Order: utils.ToJsonString(b.takerOrder),
	}

	transaction, launchLog := s.marketHandler.handleNewOrder(&takerOrderEvent)

	newTradesCount := models.TradeDao.Count()
	newTransactionsCount := models.TransactionDao.Count()

	s.Equal(b.expectedTradesCount, newTradesCount-oldTradesCount)
	s.Equal(b.expectedTransactionsCount, newTransactionsCount-oldTransactionsCount)

	if b.whenPending != nil {
		s.assertExpectedResult(b, b.whenPending)
	}

	return transaction, launchLog
}

func (s *marketHandlerSuite) TestMatchOrders0() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "100"},
		[]*buildOrderParams{
			{"buy", "140", "140"},
		},
		1,
		1,
		&expectedResult{
			[][]string{
				{"0", "100", "0", "0"},
				{"40", "100", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"140",
					"140",
				},
				{
					"buy",
					2,
					"140",
					"-100",
				},
			},
		},
		nil,
		nil,
	)
}

// 1 v 1
// taker full filled
// maker partial filled
func (s *marketHandlerSuite) TestMatchOrders1() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "100"},
		[]*buildOrderParams{
			{"buy", "140", "140"},
		},
		1,
		1,
		&expectedResult{
			[][]string{
				{"0", "100", "0", "0"},
				{"40", "100", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"140",
					"140",
				},
				{
					"buy",
					2,
					"140",
					"-100",
				},
			},
		},
		&expectedResult{
			[][]string{
				{"0", "0", "100", "0"},
				{"40", "0", "100", "0"},
			},
			[]string{common.ORDER_FULL_FILLED, common.ORDER_PENDING},
			nil,
		},
		&expectedResult{
			[][]string{
				{"0", "0", "0", "100"},
				{"40", "0", "0", "100"},
			},
			[]string{common.ORDER_CANCELED, common.ORDER_PENDING},
			nil,
		},
	)
}

// 1 v 1
// taker full filled
// maker full filled
func (s *marketHandlerSuite) TestMatchOrders2() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "80"},
		[]*buildOrderParams{
			&buildOrderParams{"buy", "141", "80"},
		},
		1,
		1,
		&expectedResult{
			[][]string{
				{"0", "80", "0", "0"},
				{"0", "80", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"141",
					"80",
				},
				{
					"buy",
					2,
					"141",
					"-80",
				},
			}},
		&expectedResult{
			[][]string{
				{"0", "0", "80", "0"},
				{"0", "0", "80", "0"},
			},
			[]string{common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED},
			nil,
		},
		&expectedResult{
			[][]string{
				{"0", "0", "0", "80"},
				{"0", "0", "0", "80"},
			},
			[]string{common.ORDER_CANCELED, common.ORDER_CANCELED},
			nil,
		},
	)
}

// 1 v 1
// taker partial filled
// maker full filled
func (s *marketHandlerSuite) TestMatchOrders3() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "100"},
		[]*buildOrderParams{
			&buildOrderParams{"buy", "141", "80"},
		},
		1,
		1,
		&expectedResult{
			[][]string{
				{"20", "80", "0", "0"},
				{"0", "80", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"141",
					"80",
				},
				{
					"buy",
					2,
					"141",
					"-80",
				},
				{
					"sell",
					3,
					"140",
					"20",
				},
			},
		},
		&expectedResult{
			[][]string{
				{"20", "0", "80", "0"},
				{"0", "0", "80", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_FULL_FILLED},
			nil,
		},
		&expectedResult{
			[][]string{
				{"20", "0", "0", "80"},
				{"0", "0", "0", "80"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_CANCELED},
			nil,
		},
	)
}

// 1 v n
// taker partial filled
// maker full filled
func (s *marketHandlerSuite) TestMatchOrders4() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "100"},
		[]*buildOrderParams{
			&buildOrderParams{"buy", "143", "20"},
			&buildOrderParams{"buy", "142", "60"},
			&buildOrderParams{"buy", "141", "10"},
		},
		3,
		1,
		&expectedResult{
			[][]string{
				{"10", "90", "0", "0"},
				{"0", "20", "0", "0"},
				{"0", "60", "0", "0"},
				{"0", "10", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"143",
					"20",
				},
				{
					"buy",
					2,
					"142",
					"60",
				},
				{
					"buy",
					3,
					"141",
					"10",
				},
				{
					"buy",
					4,
					"143",
					"-20",
				},
				{
					"buy",
					5,
					"142",
					"-60",
				},
				{
					"buy",
					6,
					"141",
					"-10",
				},
				{
					"sell",
					7,
					"140",
					"10",
				},
			},
		},
		&expectedResult{
			[][]string{
				{"10", "0", "90", "0"},
				{"0", "0", "20", "0"},
				{"0", "0", "60", "0"},
				{"0", "0", "10", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED},
			nil,
		},
		&expectedResult{
			[][]string{
				{"10", "0", "0", "90"},
				{"0", "0", "0", "20"},
				{"0", "0", "0", "60"},
				{"0", "0", "0", "10"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_CANCELED, common.ORDER_CANCELED, common.ORDER_CANCELED},
			nil,
		},
	)
}

// 1 v n
// taker full filled
// maker full filled
func (s *marketHandlerSuite) TestMatchOrders5() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "100"},
		[]*buildOrderParams{
			&buildOrderParams{"buy", "143", "20"},
			&buildOrderParams{"buy", "142", "60"},
			&buildOrderParams{"buy", "141", "20"},
		},
		3,
		1,
		&expectedResult{
			[][]string{
				{"0", "100", "0", "0"},
				{"0", "20", "0", "0"},
				{"0", "60", "0", "0"},
				{"0", "20", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"143",
					"20",
				},
				{
					"buy",
					2,
					"142",
					"60",
				},
				{
					"buy",
					3,
					"141",
					"20",
				},
				{
					"buy",
					4,
					"143",
					"-20",
				},
				{
					"buy",
					5,
					"142",
					"-60",
				},
				{
					"buy",
					6,
					"141",
					"-20",
				},
			},
		},
		&expectedResult{
			[][]string{
				{"0", "0", "100", "0"},
				{"0", "0", "20", "0"},
				{"0", "0", "60", "0"},
				{"0", "0", "20", "0"},
			},
			[]string{common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED},
			nil,
		},
		&expectedResult{
			[][]string{
				{"0", "0", "0", "100"},
				{"0", "0", "0", "20"},
				{"0", "0", "0", "60"},
				{"0", "0", "0", "20"},
			},
			[]string{common.ORDER_CANCELED, common.ORDER_CANCELED, common.ORDER_CANCELED, common.ORDER_CANCELED},
			nil,
		},
	)
}

// 1 v n
// taker full filled
// maker partial filled
func (s *marketHandlerSuite) TestMatchOrders6() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "140", "50"},
		[]*buildOrderParams{
			&buildOrderParams{"buy", "143", "20"},
			&buildOrderParams{"buy", "142", "60"},
			&buildOrderParams{"buy", "141", "20"},
		},
		2,
		1,
		&expectedResult{
			[][]string{
				{"0", "50", "0", "0"},
				{"0", "20", "0", "0"},
				{"30", "30", "0", "0"},
				{"20", "0", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING},

			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"143",
					"20",
				},
				{
					"buy",
					2,
					"142",
					"60",
				},
				{
					"buy",
					3,
					"141",
					"20",
				},
				{
					"buy",
					4,
					"143",
					"-20",
				},
				{
					"buy",
					5,
					"142",
					"-30",
				},
			},
		},
		&expectedResult{
			[][]string{
				{"0", "0", "50", "0"},
				{"0", "0", "20", "0"},
				{"30", "0", "30", "0"},
				{"20", "0", "0", "0"},
			},
			[]string{common.ORDER_FULL_FILLED, common.ORDER_FULL_FILLED, common.ORDER_PENDING, common.ORDER_PENDING},
			nil,
		},
		&expectedResult{
			[][]string{
				{"0", "0", "0", "50"},
				{"0", "0", "0", "20"},
				{"30", "0", "0", "30"},
				{"20", "0", "0", "0"},
			},
			[]string{common.ORDER_CANCELED, common.ORDER_CANCELED, common.ORDER_PENDING, common.ORDER_PENDING},
			nil,
		},
	)
}

// no match
func (s *marketHandlerSuite) TestMatchOrders7() {
	s.newBatchMatchOrdersTest(
		&buildOrderParams{"sell", "150", "50"},
		[]*buildOrderParams{
			&buildOrderParams{"buy", "143", "20"},
			&buildOrderParams{"buy", "142", "60"},
			&buildOrderParams{"buy", "141", "20"},
		},
		0,
		0,
		&expectedResult{
			[][]string{
				{"50", "0", "0", "0"},
				{"20", "0", "0", "0"},
				{"60", "0", "0", "0"},
				{"20", "0", "0", "0"},
			},
			[]string{common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING, common.ORDER_PENDING},
			[]*common.WebsocketMarketOrderChangePayload{
				{
					"buy",
					1,
					"143",
					"20",
				},
				{
					"buy",
					2,
					"142",
					"60",
				},
				{
					"buy",
					3,
					"141",
					"20",
				},
				{
					"sell",
					4,
					"150",
					"50",
				},
			},
		},
		nil, nil,
	)
}

func (s *marketHandlerSuite) assertOrderAmounts(available, pending, confirmed, canceled string, order *models.Order) {
	s.Equal(available, order.AvailableAmount.String(), "Available Amount not match")
	s.Equal(pending, order.PendingAmount.String(), "Pending Amount not match")
	s.Equal(confirmed, order.ConfirmedAmount.String(), "Confirmed Amount not match")
	s.Equal(canceled, order.CanceledAmount.String(), "Canceled Amount not match")
}

func (s *marketHandlerSuite) TestCancelOrder() {
	order1 := newModelOrder("buy", utils.StringToDecimal("0.02"), utils.StringToDecimal("10"))
	_ = models.OrderDao.InsertOrder(order1)

	newOrderEvent := common.NewOrderEvent{
		Event: common.Event{
			Type:     common.EventNewOrder,
			MarketID: order1.MarketID,
		},
		Order: utils.ToJsonString(order1),
	}
	s.marketHandler.handleNewOrder(&newOrderEvent)
	//s.EqualValues(`{"sequence":0,"bids":[["0.02","10"]],"asks":[]}`, utils.ToJsonString(s.marketHandler.orderbook.SnapshotV2()))

	cancelOrderEvent := common.CancelOrderEvent{
		Event: common.Event{
			Type:     common.EventNewOrder,
			MarketID: order1.MarketID,
		},
		ID: order1.ID,
	}
	_, _ = s.marketHandler.handleCancelOrder(&cancelOrderEvent)
	//s.EqualValues(`{"sequence":0,"bids":[],"asks":[]}`, utils.ToJsonString(s.marketHandler.orderbook.SnapshotV2()))
	return
}

func newModelOrder(side string, price, amount decimal.Decimal) *models.Order {
	var trader string
	if side == "buy" {
		trader = fakeAccount1
	} else {
		trader = fakeAccount2
	}

	market := models.MarketHotDai()

	gasFeeInQuoteToken := decimal.Zero
	quoteTokenHugeAmount := price.Mul(amount).Mul(decimal.New(1, int32(market.QuoteTokenDecimals)))
	baseTokenHugeAmount := amount.Mul(decimal.New(1, int32(market.BaseTokenDecimals)))
	orderJson := models.OrderJSON{
		Trader:                  trader,
		Relayer:                 os.Getenv("HSK_RELAYER_ADDRESS"),
		BaseCurrency:            market.BaseTokenAddress,
		QuoteCurrency:           market.QuoteTokenAddress,
		BaseCurrencyHugeAmount:  baseTokenHugeAmount,
		QuoteCurrencyHugeAmount: quoteTokenHugeAmount,
		GasTokenHugeAmount:      gasFeeInQuoteToken,
		Data:                    ethereum.GetOrderData(1, side == "sell", false, 999999999999, 0, 0, 0, rand.Uint64(), false),
	}

	id := getHydroOrderHashHexFromOrderJson(&orderJson)
	idBytes, _ := hex.DecodeString(id)

	signature, _ := ethereum.PersonalSign(idBytes, User1PrivateKey)

	orderJson.Signature = "0x" + hex.EncodeToString(signature)

	return &models.Order{
		ID:              id,
		TraderAddress:   trader,
		MarketID:        market.ID,
		Side:            side,
		Price:           price,
		Amount:          amount,
		Status:          common.ORDER_PENDING,
		Type:            "limit",
		Version:         "hydro-v1",
		AvailableAmount: amount,
		ConfirmedAmount: decimal.Zero,
		CanceledAmount:  decimal.Zero,
		PendingAmount:   decimal.Zero,
		MakerFeeRate:    decimal.Zero,
		TakerFeeRate:    decimal.Zero,
		MakerRebateRate: decimal.Zero,
		GasFeeAmount:    decimal.Zero,
		JSON:            utils.ToJsonString(&orderJson),
		CreatedAt:       time.Now().UTC(),
	}
}

func TestMarketHandler(t *testing.T) {
	suite.Run(t, new(marketHandlerSuite))
}
