package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestFindAccountMarketTrades(t *testing.T) {
	test.PreTest()
	InitTestDB()

	_, trades := TradeDao.FindAllTrades("WETH-DAI")
	assert.EqualValues(t, 0, len(trades))

	trade := NewTrade("WETH-DAI", true)
	_ = TradeDao.InsertTrade(trade)

	var count int64
	count, trades = TradeDao.FindAccountMarketTrades(trade.Maker, "WETH-DAI", common.STATUS_SUCCESSFUL, 1, 0)
	assert.EqualValues(t, 1, len(trades))
	assert.EqualValues(t, 1, count)

	count, trades = TradeDao.FindAccountMarketTrades(trade.Maker, "WETH-DAI", common.STATUS_SUCCESSFUL, 2, 0)
	assert.EqualValues(t, 1, len(trades))
	assert.EqualValues(t, 1, count)

	//count, trades = TradeDao.FindAccountMarketTrades(trade.Maker, "WETH-DAI", common.STATUS_FAILED, 2, 0)
	//assert.EqualValues(t, 0, len(trades))
	//assert.EqualValues(t, 0, count)
}

func TestTradeDao_FindAllTrades(t *testing.T) {
	test.PreTest()
	InitTestDB()

	_, trades := TradeDao.FindAllTrades("WETH-DAI")
	assert.EqualValues(t, 0, len(trades))

	trade1 := NewTrade("WETH-DAI", true)
	trade2 := NewTrade("WETH-DAI", true)
	trade3 := NewTrade("WETH-DAI", true)
	_ = TradeDao.InsertTrade(trade1)
	_ = TradeDao.InsertTrade(trade2)
	_ = TradeDao.InsertTrade(trade3)

	_, trades = TradeDao.FindAllTrades("WETH-DAI")
	assert.EqualValues(t, 3, len(trades))
}

func TestTradeDao_InsertAndFindOneAndUpdateTrade(t *testing.T) {
	test.PreTest()
	InitTestDB()

	trade := NewTrade("WETH-DAI", true)
	err := TradeDao.InsertTrade(trade)
	assert.Nil(t, err)

	dbTrades := TradeDao.FindTradesByHash(trade.TransactionHash)
	assert.EqualValues(t, 1, len(dbTrades))
	assert.EqualValues(t, trade.TransactionHash, dbTrades[0].TransactionHash)
	assert.EqualValues(t, trade.MarketID, dbTrades[0].MarketID)
	assert.EqualValues(t, trade.Status, dbTrades[0].Status)
	assert.EqualValues(t, trade.TakerOrderID, dbTrades[0].TakerOrderID)
	assert.EqualValues(t, trade.MakerOrderID, dbTrades[0].MakerOrderID)
	assert.EqualValues(t, trade.Amount.String(), dbTrades[0].Amount.String())
	assert.EqualValues(t, trade.Price.String(), dbTrades[0].Price.String())
}

func TestTradeDao_UpdateTradeAndUpdateByStatus(t *testing.T) {
	test.PreTest()
	InitTestDB()

	trade := NewTrade("WETH-DAI", true)
	err := TradeDao.InsertTrade(trade)

	assert.Nil(t, err)

	dbTrade := TradeDao.FindTradeByID(trade.ID)
	assert.EqualValues(t, trade.ID, dbTrade.ID)
	assert.EqualValues(t, trade.Status, dbTrade.Status)

	dbTrade.Status = common.STATUS_FAILED
	_ = TradeDao.UpdateTrade(dbTrade)
	dbTrade2 := TradeDao.FindTradeByID(trade.ID)
	assert.EqualValues(t, dbTrade.ID, dbTrade2.ID)
	assert.EqualValues(t, dbTrade.Status, dbTrade2.Status)
}

func TestTradeDao_FindTradesByMarket(t *testing.T) {
	test.PreTest()
	InitTestDB()

	_, trades := TradeDao.FindAllTrades("WETH-DAI")
	assert.EqualValues(t, 0, len(trades))

	//"2006-01-02T15:04:05Z07:00"
	time1, _ := time.Parse(time.RFC3339, "2019-02-02T00:00:00Z")
	time2, _ := time.Parse(time.RFC3339, "2019-02-03T00:00:00Z")
	time3, _ := time.Parse(time.RFC3339, "2019-02-04T00:00:00Z")

	time4, _ := time.Parse(time.RFC3339, "2019-02-01T00:00:00Z")
	time5, _ := time.Parse(time.RFC3339, "2019-02-05T00:00:00Z")

	trade1 := NewTradeWithTime("WETH-DAI", true, time1)
	trade2 := NewTradeWithTime("WETH-DAI", true, time2)
	trade3 := NewTradeWithTime("WETH-DAI", true, time3)
	_ = TradeDao.InsertTrade(trade1)
	_ = TradeDao.InsertTrade(trade2)
	_ = TradeDao.InsertTrade(trade3)

	_, trades1 := TradeDao.FindAllTrades("WETH-DAI")
	assert.EqualValues(t, 3, len(trades1))

	trades2 := TradeDao.FindTradesByMarket("WETH-DAI", time4, time5)
	assert.EqualValues(t, 3, len(trades2))

	trades3 := TradeDao.FindTradesByMarket("WETH-DAI", time2, time5)
	assert.EqualValues(t, 1, len(trades3))
}

func NewTradeWithTime(marketID string, success bool, time time.Time) *Trade {
	status := common.STATUS_SUCCESSFUL

	if !success {
		status = common.STATUS_FAILED
	}

	transaction := RandomTransaction(success)
	makerOrder, takerOrder := RandomMatchOrder()

	trade := Trade{
		ID:              rand.Int63(),
		TransactionHash: transaction.TransactionHash.String,
		Status:          status,
		MarketID:        marketID,
		Maker:           makerOrder.TraderAddress,
		Taker:           takerOrder.TraderAddress,
		TakerSide:       takerOrder.Side,
		MakerOrderID:    makerOrder.ID,
		TakerOrderID:    takerOrder.ID,
		Sequence:        0,
		Amount:          takerOrder.Amount,
		Price:           takerOrder.Price,
		ExecutedAt:      time,
		CreatedAt:       time,
	}

	return &trade
}

func NewTrade(marketID string, success bool) *Trade {
	return NewTradeWithTime(marketID, success, time.Now())
}
