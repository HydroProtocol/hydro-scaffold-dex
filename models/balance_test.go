package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalanceDao_GetByAccountAndSymbol2(t *testing.T) {
	test.PreTest()
	InitTestDB()

	account := "0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"

	order1 := NewOrder(account, "WETH-DAI", "buy", false)
	_ = OrderDao.InsertOrder(order1)
	actBalance := BalanceDao.GetByAccountAndSymbol(account, "DAI", 18)
	exceptBalance := order1.AvailableAmount.Add(order1.PendingAmount).Mul(order1.Price).Mul(decimal.New(1, 18))
	assert.EqualValues(t, exceptBalance, actBalance)

	order2 := NewOrder(account, "WETH-DAI", "buy", true)
	_ = OrderDao.InsertOrder(order2)
	actBalance2 := BalanceDao.GetByAccountAndSymbol(account, "DAI", 18)
	exceptBalance2 := exceptBalance.Add(order2.PendingAmount.Add(order2.AvailableAmount).Mul(order2.Price).Mul(decimal.New(1, 18)))
	assert.EqualValues(t, exceptBalance2, actBalance2)

	order3 := NewOrder(account, "DAI-HOT", "sell", true)
	_ = OrderDao.InsertOrder(order3)
	actBalance3 := BalanceDao.GetByAccountAndSymbol(account, "DAI", 18)
	exceptBalance3 := exceptBalance2.Add(order3.AvailableAmount.Add(order3.PendingAmount).Mul(decimal.New(1, 18)))
	assert.EqualValues(t, exceptBalance3, actBalance3)
}
