package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	uuid2 "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestGetAccountOrders(t *testing.T) {
	test.PreTest()
	InitTestDB()

	orders := OrderDao.FindMarketPendingOrders("WETH-DAI")
	assert.EqualValues(t, 0, len(orders))

	order1 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order2 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order3 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order4 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order5 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order6 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order7 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order8 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order9 := NewOrder(config.User1, "WETH-DAI", "buy", false)

	order10 := NewOrder(config.User2, "WETH-DAI", "buy", false)

	_ = OrderDao.InsertOrder(order1)
	_ = OrderDao.InsertOrder(order2)
	_ = OrderDao.InsertOrder(order3)
	_ = OrderDao.InsertOrder(order4)
	_ = OrderDao.InsertOrder(order5)
	_ = OrderDao.InsertOrder(order6)
	_ = OrderDao.InsertOrder(order7)
	_ = OrderDao.InsertOrder(order8)
	_ = OrderDao.InsertOrder(order9)
	_ = OrderDao.InsertOrder(order10)

	var count int64
	count, orders = OrderDao.FindByAccount(config.User1, "WETH-DAI", common.ORDER_PENDING, 9, 3)
	assert.EqualValues(t, 6, len(orders))
	assert.EqualValues(t, 9, count)

	count, orders = OrderDao.FindByAccount(config.User1, "WETH-DAI", common.ORDER_PENDING, 10, 0)
	assert.EqualValues(t, 9, len(orders))
	assert.EqualValues(t, 9, count)

	count, orders = OrderDao.FindByAccount(config.User1, "WETH-DAI", common.ORDER_PENDING, 9, 0)
	assert.EqualValues(t, 9, len(orders))
	assert.EqualValues(t, 9, count)

	count, orders = OrderDao.FindByAccount(config.User1, "WETH-DAI", common.ORDER_FULL_FILLED, 9, 0)
	assert.EqualValues(t, 0, len(orders))
	assert.EqualValues(t, 0, count)
}

func TestGetMarketPendingOrders(t *testing.T) {
	test.PreTest()
	InitTestDB()

	orders := OrderDao.FindMarketPendingOrders("WETH-DAI")
	assert.EqualValues(t, 0, len(orders))

	order1 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order2 := NewOrder(config.User1, "WETH-DAI", "buy", false)
	order3 := NewOrder(config.User1, "WETH-DAI", "buy", false)

	_ = OrderDao.InsertOrder(order1)
	_ = OrderDao.InsertOrder(order2)
	_ = OrderDao.InsertOrder(order3)

	orders = OrderDao.FindMarketPendingOrders("WETH-DAI")
	assert.EqualValues(t, 3, len(orders))
}

func TestFindNotExistOrder(t *testing.T) {
	test.PreTest()
	InitTestDB()

	dbOrder := OrderDao.FindByID("empty_id")
	assert.Nil(t, dbOrder)

}

func TestInsertAndFindOneAndUpdateOrders(t *testing.T) {
	test.PreTest()
	InitTestDB()

	order := RandomOrder()

	err := OrderDao.InsertOrder(order)
	assert.Nil(t, err)

	dbOrder := OrderDao.FindByID(order.ID)
	assert.EqualValues(t, dbOrder.ID, order.ID)
	assert.EqualValues(t, dbOrder.Status, order.Status)
	assert.EqualValues(t, dbOrder.Amount.String(), order.Amount.String())
	assert.EqualValues(t, dbOrder.Price.String(), order.Price.String())
	assert.EqualValues(t, dbOrder.AvailableAmount.String(), order.AvailableAmount.String())
	assert.EqualValues(t, dbOrder.PendingAmount.String(), order.PendingAmount.String())

	dbOrder.PendingAmount.Add(dbOrder.AvailableAmount)
	dbOrder.AvailableAmount = decimal.Zero
	err = OrderDao.UpdateOrder(dbOrder)
	dbOrder2 := OrderDao.FindByID(order.ID)

	assert.EqualValues(t, dbOrder.AvailableAmount.String(), dbOrder2.AvailableAmount.String())
	assert.EqualValues(t, dbOrder.PendingAmount.String(), dbOrder2.PendingAmount.String())
}

func TestOrder_GetOrderJson(t *testing.T) {
	json := OrderJSON{
		Trader:                  config.User1,
		Relayer:                 config.Getenv("HSK_RELAYER_ADDRESS"),
		BaseCurrencyHugeAmount:  utils.StringToDecimal("100000000000000000000000000000000000"),
		QuoteCurrencyHugeAmount: utils.StringToDecimal("200000000000000000000000000000000000"),
		BaseCurrency:            config.Getenv("HSK_HYDRO_TOKEN_ADDRESS"),
		QuoteCurrency:           config.Getenv("HSK_USD_TOKEN_ADDRESS"),
		GasTokenHugeAmount:      utils.StringToDecimal("1000000000"),
		Signature:               "0x15a85430057580a5a35125db098b686b3541a291b3fce69365dc47d502fa63395ce9f7100240e4558c6ad29b8aa9a2c01d2b5353babdffd6ac50babf0127fdd600",
		Data:                    "something",
	}
	jsonStr := utils.ToJsonString(json)

	order := RandomOrder()
	order.JSON = jsonStr

	assert.EqualValues(t, json.Trader, order.GetOrderJson().Trader)
	assert.EqualValues(t, json.Relayer, order.GetOrderJson().Relayer)
	assert.EqualValues(t, json.BaseCurrencyHugeAmount, order.GetOrderJson().BaseCurrencyHugeAmount)
	assert.EqualValues(t, json.QuoteCurrencyHugeAmount, order.GetOrderJson().QuoteCurrencyHugeAmount)
	assert.EqualValues(t, json.Signature, order.GetOrderJson().Signature)
}

func NewOrder(account, marketID, side string, withPending bool) *Order {
	id := uuid2.NewV4().String()
	amountInt := rand.Intn(10) + 2
	amount := utils.IntToDecimal(amountInt)

	pendingAmountInt := 0
	if withPending {
		pendingAmountInt = rand.Intn(amountInt-1) + 1
	}

	order := &Order{
		ID:              id,
		TraderAddress:   account,
		MarketID:        marketID,
		Side:            side,
		Type:            "limit",
		Price:           utils.IntToDecimal(rand.Intn(100) + 50),
		Amount:          amount,
		Status:          common.ORDER_PENDING,
		Version:         "v1",
		AvailableAmount: utils.IntToDecimal(amountInt - pendingAmountInt),
		ConfirmedAmount: utils.StringToDecimal("0"),
		CanceledAmount:  utils.StringToDecimal("0"),
		PendingAmount:   utils.IntToDecimal(pendingAmountInt),
		TakerFeeRate:    utils.StringToDecimal("0.003"),
		MakerFeeRate:    utils.StringToDecimal("0.001"),
		MakerRebateRate: utils.StringToDecimal("0"),
		GasFeeAmount:    utils.StringToDecimal("1000000"),
		JSON:            "something",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return order
}

func RandomMatchOrder() (*Order, *Order) {
	makerOrder := RandomOrder()
	side := "buy"
	if makerOrder.Side == "buy" {
		side = "sell"
	}

	takerOrder := &Order{
		ID:              uuid2.NewV4().String(),
		TraderAddress:   makerOrder.TraderAddress,
		MarketID:        makerOrder.MarketID,
		Side:            side,
		Type:            "limit",
		Price:           makerOrder.Price,
		Amount:          makerOrder.Amount,
		Status:          common.ORDER_PENDING,
		Version:         "v1",
		AvailableAmount: makerOrder.Amount,
		ConfirmedAmount: utils.StringToDecimal("0"),
		CanceledAmount:  utils.StringToDecimal("0"),
		PendingAmount:   utils.StringToDecimal("0"),
		TakerFeeRate:    utils.StringToDecimal("0.003"),
		MakerFeeRate:    utils.StringToDecimal("0.001"),
		MakerRebateRate: utils.StringToDecimal("0"),
		GasFeeAmount:    utils.StringToDecimal("1000000"),
		JSON:            "something",
		CreatedAt:       time.Now(),
	}

	return makerOrder, takerOrder
}

func RandomOrder() *Order {
	markets := []string{"WETH-DAI", "HOT-DAI", "AIR-DAI", "DAI-WETH", "HOT-WETH", "AIR-WETH", "TRX-DAI", "TRX-WETH"}
	accounts := []string{"0xe36ea790bc9d7ab70c55260c66d52b1eca985f84", "0xe834ec434daba538cd1b9fe1582052b880bd7e63", "0x78dc5d2d739606d31509c31d654056a45185ecb6", "0xa8dda8d7f5310e4a9e24f8eba77e091ac264f872", "0x06cef8e666768cc40cc78cf93d9611019ddcb628", "0x4404ac8bd8f9618d27ad2f1485aa1b2cfd82482d", "0x7457d5e02197480db681d3fdf256c7aca21bdc12"}
	sides := []string{"buy", "sell"}
	types := []string{"limit", "market"}

	id := uuid2.NewV4().String()
	amount := utils.IntToDecimal(rand.Intn(10) + 1)
	order := &Order{
		ID:              id,
		TraderAddress:   accounts[rand.Intn(len(accounts))],
		MarketID:        markets[rand.Intn(len(markets))],
		Side:            sides[rand.Intn(len(sides))],
		Type:            types[rand.Intn(len(types))],
		Price:           utils.IntToDecimal(rand.Intn(100) + 50),
		Amount:          amount,
		Status:          common.ORDER_PENDING,
		Version:         "v1",
		AvailableAmount: amount,
		ConfirmedAmount: utils.StringToDecimal("0"),
		CanceledAmount:  utils.StringToDecimal("0"),
		PendingAmount:   utils.StringToDecimal("0"),
		TakerFeeRate:    utils.StringToDecimal("0.003"),
		MakerFeeRate:    utils.StringToDecimal("0.001"),
		MakerRebateRate: utils.StringToDecimal("0"),
		GasFeeAmount:    utils.StringToDecimal("1000000"),
		JSON:            "something",
		CreatedAt:       time.Now(),
	}

	return order
}
