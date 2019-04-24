package models

import (
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBalanceDao_PG_GetByAccountAndSymbol(t *testing.T) {
	setEnvs()
	InitTestDBPG()

	account := "0xe36ea790bc9d7ab70c55260c66d52b1eca985f84"

	order1 := NewOrder(account, "WETH-DAI", "buy", false)
	_ = OrderDaoPG.InsertOrder(order1)
	actBalance := BalanceDaoPG.GetByAccountAndSymbol(account, "DAI", 18)
	exceptBalance := order1.AvailableAmount.Add(order1.PendingAmount).Mul(order1.Price).Mul(decimal.New(1, 18))
	assert.EqualValues(t, exceptBalance.String(), actBalance.String())

	order2 := NewOrder(account, "WETH-DAI", "buy", true)
	_ = OrderDaoPG.InsertOrder(order2)
	actBalance2 := BalanceDaoPG.GetByAccountAndSymbol(account, "DAI", 18)
	exceptBalance2 := exceptBalance.Add(order2.PendingAmount.Add(order2.AvailableAmount).Mul(order2.Price).Mul(decimal.New(1, 18)))
	assert.EqualValues(t, exceptBalance2.String(), actBalance2.String())

	order3 := NewOrder(account, "DAI-HOT", "sell", true)
	_ = OrderDaoPG.InsertOrder(order3)
	actBalance3 := BalanceDaoPG.GetByAccountAndSymbol(account, "DAI", 18)
	exceptBalance3 := exceptBalance2.Add(order3.AvailableAmount.Add(order3.PendingAmount).Mul(decimal.New(1, 18)))
	assert.EqualValues(t, exceptBalance3.String(), actBalance3.String())
}
