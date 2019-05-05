package dex_engine

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
)

func UpdateOrder(order *models.Order) error {
	err := models.OrderDao.UpdateOrder(order)
	market := models.MarketDao.FindMarketByID(order.MarketID)
	sendOrderUpdateMessage(order)
	if order.Side == "buy" {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.QuoteTokenSymbol, models.BalanceDao.GetByAccountAndSymbol(order.TraderAddress, market.QuoteTokenSymbol, market.QuoteTokenDecimals))
	} else {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.BaseTokenSymbol, models.BalanceDao.GetByAccountAndSymbol(order.TraderAddress, market.BaseTokenSymbol, market.BaseTokenDecimals))
	}

	return err
}

func InsertOrder(order *models.Order) error {
	err := models.OrderDao.InsertOrder(order)
	market := models.MarketDao.FindMarketByID(order.MarketID)
	sendOrderUpdateMessage(order)

	if order.Side == "buy" {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.QuoteTokenSymbol, models.BalanceDao.GetByAccountAndSymbol(order.TraderAddress, market.QuoteTokenSymbol, market.QuoteTokenDecimals))
	} else {
		sendLockedBalanceChangeMessage(order.TraderAddress, market.BaseTokenSymbol, models.BalanceDao.GetByAccountAndSymbol(order.TraderAddress, market.BaseTokenSymbol, market.BaseTokenDecimals))
	}

	return err
}

func UpdateTrade(trade *models.Trade) error {
	err := models.TradeDao.UpdateTrade(trade)
	sendTradeUpdateMessage(trade)

	if trade.Status == common.STATUS_SUCCESSFUL {
		sendNewMarketTradeMessage(trade)
	}
	return err
}

func InsertTrade(trade *models.Trade) error {
	err := models.TradeDao.InsertTrade(trade)
	sendTradeUpdateMessage(trade)
	return err
}

type MatchResultWithOrders struct {
	Sequence uint64
	*common.MatchResult
	modelTakerOrder  *models.Order
	modelMakerOrders map[string]*models.Order
}

func NewMatchResultWithOrders(takerOrder *models.Order, result *common.MatchResult) *MatchResultWithOrders {
	r := &MatchResultWithOrders{}

	r.MatchResult = result
	r.modelTakerOrder = takerOrder
	r.modelMakerOrders = make(map[string]*models.Order)

	for i := range result.MatchItems {
		item := result.MatchItems[i]
		r.modelMakerOrders[item.MakerOrder.ID] = models.OrderDao.FindByID(item.MakerOrder.ID)
	}

	return r
}
