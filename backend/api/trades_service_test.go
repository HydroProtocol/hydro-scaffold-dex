package api

import (
	"fmt"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBuildTradingViewByTrades(t *testing.T) {
	var trades []*models.Trade
	trades = append(trades, newTestTrade("0.1", "1", 100))
	trades = append(trades, newTestTrade("0.2", "1", 1002))
	trades = append(trades, newTestTrade("0.3", "1", 1000))
	trades = append(trades, newTestTrade("0.01", "1", 1000))
	trades = append(trades, newTestTrade("0.5", "1", 1000))
	trades = append(trades, newTestTrade("0.6", "1", 1000))
	trades = append(trades, newTestTrade("0.3", "1", 1005))
	trades = append(trades, newTestTrade("0.3", "1", 3005))

	bars := BuildTradingViewByTrades(trades, 3000)
	assert.EqualValues(t, `[{"time":0,"open":"0.1","close":"0.3","low":"0.01","high":"0.6","volume":"7"},{"time":3000,"open":"0.3","close":"0.3","low":"0.3","high":"0.3","volume":"1"}]`, utils.ToJsonString(bars))

	trades = append(trades, newTestTrade("0.3", "1", 9005))
	bars2 := BuildTradingViewByTrades(trades, 3000)
	assert.EqualValues(t, `[{"time":0,"open":"0.1","close":"0.3","low":"0.01","high":"0.6","volume":"7"},{"time":3000,"open":"0.3","close":"0.3","low":"0.3","high":"0.3","volume":"1"},{"time":9000,"open":"0.3","close":"0.3","low":"0.3","high":"0.3","volume":"1"}]`, utils.ToJsonString(bars2))

	fmt.Println(utils.ToJsonString(trades))
}

func newTestTrade(price, amount string, executedAt int64) *models.Trade {
	priceDecimal, _ := decimal.NewFromString(price)
	amountDecimal, _ := decimal.NewFromString(amount)
	executedAtTime := time.Unix(executedAt, 0)
	return &models.Trade{
		Price:      priceDecimal,
		Amount:     amountDecimal,
		ExecutedAt: executedAtTime,
	}
}
