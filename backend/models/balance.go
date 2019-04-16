package models

import (
	"github.com/shopspring/decimal"
)

// IBalanceDao is an interface about how to fetch balance data from storage
type IBalanceDao interface {
	GetByAccountAndSymbol(account, tokenSymbol string, decimals int) decimal.Decimal
}

// balanceDao is default dao to fetch balance data from db.
type balanceDao struct {
	IBalanceDao
}

var BalanceDao IBalanceDao

func init() {
	BalanceDao = &balanceDao{}
}

type nullDecimal struct {
	value decimal.Decimal
}

func (d *nullDecimal) Scan(value interface{}) error {
	if value == nil {
		d.value = decimal.Zero
		return nil
	} else {
		return d.value.Scan(value)
	}
}

func (_ *balanceDao) GetByAccountAndSymbol(account, tokenSymbol string, decimals int) decimal.Decimal {
	var sellLockedBalance nullDecimal
	var buyLockedBalance nullDecimal

	err := DB.QueryRow(`select sum(amount) from orders where status='pending' and trader_address= $1 and market_id like $2 and side = 'sell'`, account, tokenSymbol+"-%").Scan(&sellLockedBalance)
	if err != nil {
		panic(err)
	}

	err = DB.QueryRow(`select sum( (available_amount + pending_amount) * price) from orders where trader_address = $1 and status = 'pending' and market_id like $2 and side = 'buy'`, account, "%-"+tokenSymbol).Scan(&buyLockedBalance)
	if err != nil {
		panic(err)
	}

	return sellLockedBalance.value.Add(buyLockedBalance.value).Mul(decimal.New(1, int32(decimals)))
}
