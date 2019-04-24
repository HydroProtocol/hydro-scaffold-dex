package models

import (
	"github.com/shopspring/decimal"
)

// IBalanceDao is an interface about how to fetch balance data from storage
type IBalanceDao interface {
	GetByAccountAndSymbol(account, tokenSymbol string, decimals int) decimal.Decimal
}

var BalanceDao IBalanceDao
var BalanceDaoPG IBalanceDao

func init() {
	BalanceDao = &balanceDaoPG{}
	BalanceDaoPG = BalanceDao
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

type balanceDaoPG struct {
}

func (balanceDaoPG) GetByAccountAndSymbol(account, tokenSymbol string, decimals int) decimal.Decimal {
	var sellLockedBalance nullDecimal
	var buyLockedBalance nullDecimal

	sellRow := DB.Raw(`select sum(amount) as locked_balance from orders where status='pending' and trader_address= $1 and market_id like $2 and side = 'sell'`, account, tokenSymbol+"-%").Row()
	if sellRow == nil {
		sellLockedBalance.Scan(nil)
	}
	err := sellRow.Scan(&sellLockedBalance)
	if err != nil {
		panic(err)
	}

	buyRow := DB.Raw(`select sum( (available_amount + pending_amount) * price) as locked_balance from orders where trader_address = $1 and status = 'pending' and market_id like $2 and side = 'buy'`, account, "%-"+tokenSymbol).Row()
	if buyRow == nil {
		buyLockedBalance.Scan(nil)
	}
	err = buyRow.Scan(&buyLockedBalance)
	if err != nil {
		panic(err)
	}

	return sellLockedBalance.value.Add(buyLockedBalance.value).Mul(decimal.New(1, int32(decimals)))
}
