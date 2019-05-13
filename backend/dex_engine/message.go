package dex_engine

import (
	"encoding/json"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
)

// This queue is used to send message to ws servers
// todo a little weird global var
var wsQueue common.IQueue = nil

func InitWsQueue(queue common.IQueue) {
	wsQueue = queue
}

func sendOrderUpdateMessage(order *models.Order) {
	_ = pushAccountMessage(order.TraderAddress, &common.WebsocketOrderChangePayload{
		Type:  common.WsTypeOrderChange,
		Order: order,
	})
}

func sendTradeUpdateMessage(trade *models.Trade) {
	_ = pushAccountMessage(trade.Maker, &common.WebsocketTradeChangePayload{
		Type:  common.WsTypeTradeChange,
		Trade: trade,
	})

	_ = pushAccountMessage(trade.Taker, &common.WebsocketTradeChangePayload{
		Type:  common.WsTypeTradeChange,
		Trade: trade,
	})
}

func sendNewMarketTradeMessage(trade *models.Trade) {
	_ = pushMarketChannel(trade.MarketID, &common.WebsocketMarketNewMarketTradePayload{
		Type:  common.WsTypeNewMarketTrade,
		Trade: trade,
	})
}

func sendLockedBalanceChangeMessage(address, symbol string, newLockedBalance decimal.Decimal) {
	_ = pushAccountMessage(address, &common.WebsocketLockedBalanceChangePayload{
		Type:    common.WsTypeLockedBalanceChange,
		Symbol:  symbol,
		Balance: newLockedBalance,
	})
}

//func sendOrderbookChangeMessage(marketID string, sequence uint64, side string, price, amount decimal.Decimal) error {
//	payload := &common.WebsocketMarketOrderChangePayload{
//		Sequence: sequence,
//		Side:     side,
//		Price:    price.String(),
//		Amount:   amount.String(),
//	}
//
//	return pushMarketChannel(marketID, payload)
//}

func pushMarketChannel(marketID string, payload interface{}) error {
	return pushMessage(&common.WebSocketMessage{
		ChannelID: common.GetMarketChannelID(marketID),
		Payload:   payload,
	})
}

func pushAccountMessage(address string, payload interface{}) error {
	return pushMessage(&common.WebSocketMessage{
		ChannelID: common.GetAccountChannelID(address),
		Payload:   payload,
	})
}

func pushMessage(message interface{}) error {
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return nil
	}

	utils.Debugf("sending pushMessage: %v", string(msgBytes))
	return wsQueue.Push(msgBytes)
}
