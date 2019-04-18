package augur_engine

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
)

//type Orderbook struct {
//	*common.Orderbook
//	Sequence uint64
//}

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

//func NewOrderbook(marketID string) *Orderbook {
//	originalOrderbook := common.NewOrderbook(marketID)
//
//	orderbook := &Orderbook{
//		Orderbook: originalOrderbook,
//		Sequence:  uint64(0),
//	}
//
//	return orderbook
//}
