package dex_engine

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
)

func UpdateOrder(order *models.Order) error {
	err := models.OrderDaoSqlite.UpdateOrder(order)
	market := models.MarketDaoSqlite.FindMarketByID(order.MarketID)
	sendOrderUpdateMessage(order)
	if order.Side == "buy" {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.QuoteTokenSymbol, models.BalanceDaoSqlite.GetByAccountAndSymbol(order.TraderAddress, market.QuoteTokenSymbol, market.QuoteTokenDecimals))
	} else {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.BaseTokenSymbol, models.BalanceDaoSqlite.GetByAccountAndSymbol(order.TraderAddress, market.BaseTokenSymbol, market.BaseTokenDecimals))
	}

	return err
}

func InsertOrder(order *models.Order) error {
	err := models.OrderDaoSqlite.InsertOrder(order)
	market := models.MarketDaoSqlite.FindMarketByID(order.MarketID)
	sendOrderUpdateMessage(order)

	if order.Side == "buy" {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.QuoteTokenSymbol, models.BalanceDaoSqlite.GetByAccountAndSymbol(order.TraderAddress, market.QuoteTokenSymbol, market.QuoteTokenDecimals))
	} else {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.BaseTokenSymbol, models.BalanceDaoSqlite.GetByAccountAndSymbol(order.TraderAddress, market.BaseTokenSymbol, market.BaseTokenDecimals))
	}

	return err
}

func UpdateTrade(trade *models.Trade) error {
	err := models.TradeDaoSqlite.UpdateTrade(trade)
	sendTradeUpdateMessage(trade)

	if trade.Status == common.STATUS_SUCCESSFUL {
		sendNewMarketTradeMessage(trade)
	}
	return err
}

func InsertTrade(trade *models.Trade) error {
	err := models.TradeDaoSqlite.InsertTrade(trade)
	sendTradeUpdateMessage(trade)
	return err
}
