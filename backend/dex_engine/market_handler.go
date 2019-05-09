package dex_engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HydroProtocol/hydro-sdk-backend/engine"
	"math/big"
	"os"
	"runtime"
	"time"

	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
)

type MarketHandler struct {
	ctx         context.Context
	market      *models.Market
	eventChan   chan []byte
	hydroEngine *engine.Engine
}

// Run is synchronous, it will be improved in the later releases.
func (m *MarketHandler) Run() {
	for data := range m.eventChan {
		_ = handleEvent(m, string(data))
	}
	utils.Infof("market %s stopped", m.market.ID)
}

func (m *MarketHandler) Stop() {
	close(m.eventChan)
}

// handleEvent recover any panic which is caused by event.
// It will log event and response as well.
func handleEvent(marketHandler *MarketHandler, eventJSON string) (err error) {
	var event common.Event

	defer func() {
		if rcv := recover(); rcv != nil {
			err = rcv.(error)
		}

		if err != nil {
			buf := make([]byte, 2048)
			n := runtime.Stack(buf, false)
			stackInfo := fmt.Sprintf("%s", buf[:n])

			utils.Errorf("Errorf: %+v", err)
			utils.Errorf(stackInfo)
		}

	}()

	err = json.Unmarshal([]byte(eventJSON), &event)

	if err != nil {
		utils.Errorf("Unmarshal event failed %s", eventJSON)
		return err
	}

	_, err = marketHandler.handleEvent(event, eventJSON)

	return err
}

func (m *MarketHandler) handleEvent(event common.Event, eventJSON string) (interface{}, error) {
	switch event.Type {
	case common.EventNewOrder:
		var e common.NewOrderEvent
		_ = json.Unmarshal([]byte(eventJSON), &e)
		res, _ := m.handleNewOrder(&e)
		return res, nil
	case common.EventCancelOrder:
		var e common.CancelOrderEvent
		_ = json.Unmarshal([]byte(eventJSON), &e)
		res, err := m.handleCancelOrder(&e)
		return res, err
	case common.EventConfirmTransaction:
		var e common.ConfirmTransactionEvent
		_ = json.Unmarshal([]byte(eventJSON), &e)
		res, err := m.handleTransactionResult(&e)
		return res, err
	default:
		return nil, fmt.Errorf("unsupport event for market %s %s", m.market.ID, eventJSON)
	}
}

func (m MarketHandler) handleNewOrder(event *common.NewOrderEvent) (transaction *models.Transaction, launchLog *models.LaunchLog) {
	eventOrderString := event.Order
	var eventOrder models.Order
	_ = json.Unmarshal([]byte(eventOrderString), &eventOrder)

	eventMemoryOrder := &common.MemoryOrder{
		ID:           eventOrder.ID,
		MarketID:     eventOrder.MarketID,
		Price:        eventOrder.Price,
		Amount:       eventOrder.Amount,
		Side:         eventOrder.Side,
		GasFeeAmount: eventOrder.GasFeeAmount,
		MakerFeeRate: eventOrder.MakerFeeRate,
		TakerFeeRate: eventOrder.TakerFeeRate,
	}

	utils.Debugf("%s NEW_ORDER  price: %s amount: %s %4s", event.MarketID, eventOrder.Price.StringFixed(5), eventOrder.Amount.StringFixed(5), eventOrder.Side)

	matchResult, hasMatch := m.hydroEngine.HandleNewOrder(eventMemoryOrder)
	if hasMatch {
		resultWithOrders := NewMatchResultWithOrders(&eventOrder, &matchResult)

		for i := range resultWithOrders.MatchItems {
			item := resultWithOrders.MatchItems[i]
			makerOrder := resultWithOrders.modelMakerOrders[item.MakerOrder.ID]

			makerOrder.AvailableAmount = makerOrder.AvailableAmount.Sub(item.MatchedAmount)
			eventOrder.AvailableAmount = eventOrder.AvailableAmount.Sub(item.MatchedAmount)

			if item.MatchShouldBeCanceled {
				makerOrder.CanceledAmount = makerOrder.CanceledAmount.Add(item.MatchedAmount)
				eventOrder.CanceledAmount = eventOrder.CanceledAmount.Add(item.MatchedAmount)
			} else {
				makerOrder.PendingAmount = makerOrder.PendingAmount.Add(item.MatchedAmount)
				eventOrder.PendingAmount = eventOrder.PendingAmount.Add(item.MatchedAmount)
			}

			if item.MakerOrderIsDone {
				makerOrder.CanceledAmount = makerOrder.Amount.Sub(makerOrder.ConfirmedAmount.Add(makerOrder.PendingAmount))
				makerOrder.AvailableAmount = decimal.Zero
			}

			_ = UpdateOrder(makerOrder)

			utils.Debugf("  [Take Liquidity] price: %s amount: %s (%s) ", item.MakerOrder.Price.StringFixed(5), item.MatchedAmount.StringFixed(5), item.MakerOrder.ID)
		}

		if matchResult.TakerOrderIsDone {
			eventOrder.CanceledAmount = eventOrder.Amount.Sub(eventOrder.ConfirmedAmount.Add(eventOrder.PendingAmount))
			eventOrder.AvailableAmount = decimal.Zero
		}

		if matchResult.ExistMatchToBeExecuted() {
			transaction, launchLog = processTransactionAndLaunchLog(resultWithOrders)
			trades := newTradesByMatchResult(resultWithOrders, transaction.ID)

			for _, trade := range trades {
				_ = InsertTrade(trade)
			}
		}
	}

	_ = InsertOrder(&eventOrder)

	return transaction, launchLog
}

// If there are many items in the match result, it can't settle them in a single transaction, since there is a gas limit of a block.
// Will separate the matches into different transactions in another  release.
func processTransactionAndLaunchLog(matchResult *MatchResultWithOrders) (*models.Transaction, *models.LaunchLog) {
	takerOrder := matchResult.modelTakerOrder
	hydroTakerOrder := getHydroOrderFromModelOrder(takerOrder.GetOrderJson())

	var hydroMakerOrders []*sdk.Order
	var baseTokenFilledAmounts []*big.Int

	market := models.MarketDao.FindMarketByID(takerOrder.MarketID)

	baseTokenDecimal := market.BaseTokenDecimals

	for _, item := range matchResult.MatchItems {
		if item.MatchShouldBeCanceled {
			//skip if match should be canceled
			continue
		}

		modelMakerOrder := matchResult.modelMakerOrders[item.MakerOrder.ID]

		hydroMakerOrder := getHydroOrderFromModelOrder(modelMakerOrder.GetOrderJson())
		hydroMakerOrders = append(hydroMakerOrders, hydroMakerOrder)

		baseTokenHugeAmt := item.MatchedAmount.Mul(decimal.New(1, int32(baseTokenDecimal))).Truncate(0)
		baseTokenFilledAmt := utils.DecimalToBigInt(baseTokenHugeAmt)
		baseTokenFilledAmounts = append(baseTokenFilledAmounts, baseTokenFilledAmt)

		_ = UpdateOrder(modelMakerOrder)
	}

	transaction := &models.Transaction{
		Status: common.STATUS_PENDING,
		TransactionHash: &sql.NullString{
			Valid:  false,
			String: "",
		},
		MarketID:   takerOrder.MarketID,
		ExecutedAt: time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
	}
	err := models.TransactionDao.InsertTransaction(transaction)

	if err != nil {
		panic(err)
	}

	launchLog := &models.LaunchLog{
		ItemType:  "hydroTrade",
		ItemID:    transaction.ID,
		Status:    "created",
		From:      os.Getenv("HSK_RELAYER_ADDRESS"),
		To:        os.Getenv("HSK_HYBRID_EXCHANGE_ADDRESS"),
		Value:     decimal.Zero,
		GasLimit:  int64(len(matchResult.MatchItems) * 250000),
		Data:      utils.Bytes2HexP(hydroProtocol.GetMatchOrderCallData(hydroTakerOrder, hydroMakerOrders, baseTokenFilledAmounts)),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err = models.LaunchLogDao.InsertLaunchLog(launchLog)

	if err != nil {
		panic(err)
	}

	return transaction, launchLog
}

func newTradesByMatchResult(matchResult *MatchResultWithOrders, transactionID int64) []*models.Trade {
	var trades []*models.Trade
	takerOrder := matchResult.modelTakerOrder

	for i, item := range matchResult.MatchItems {
		modelMakerOrder := matchResult.modelMakerOrders[item.MakerOrder.ID]
		trade := &models.Trade{
			TransactionID:   transactionID,
			TransactionHash: "",
			Status:          common.STATUS_PENDING,
			MarketID:        takerOrder.MarketID,
			Maker:           modelMakerOrder.TraderAddress,
			Taker:           takerOrder.TraderAddress,
			TakerSide:       takerOrder.Side,
			MakerOrderID:    modelMakerOrder.ID,
			TakerOrderID:    takerOrder.ID,
			Sequence:        i,
			Amount:          item.MatchedAmount,
			Price:           modelMakerOrder.Price,
			CreatedAt:       time.Now().UTC(),
		}
		trades = append(trades, trade)
	}

	return trades
}

func (m *MarketHandler) handleCancelOrder(event *common.CancelOrderEvent) (interface{}, error) {
	order := models.OrderDao.FindByID(event.ID)
	if order == nil {
		return nil, errors.New(fmt.Sprintf("cannot find order with id %s", event.ID))
	}

	bookOrder := &common.MemoryOrder{
		MarketID: m.market.ID,
		ID:       order.ID,
		Price:    order.Price,
		Side:     order.Side,
		Amount:   order.AvailableAmount,
	}
	msg, success := m.hydroEngine.HandleCancelOrder(bookOrder)
	if success {
		pushMessage(msg)
	}

	order.CanceledAmount = order.CanceledAmount.Add(order.AvailableAmount)
	order.AvailableAmount = decimal.Zero
	order.AutoSetStatusByAmounts()

	err := UpdateOrder(order)

	return order, err
}

func (m *MarketHandler) handleTransactionResult(event *common.ConfirmTransactionEvent) (interface{}, error) {
	executedAt := time.Unix(int64(event.Timestamp), 0)
	transaction := models.TransactionDao.FindTransactionByHash(event.Hash)
	transaction.Status = event.Status
	transaction.ExecutedAt = executedAt
	_ = models.TransactionDao.UpdateTransaction(transaction)

	_ = models.LaunchLogDao.UpdateLaunchLogsStatusByItemID(event.Status, transaction.ID)

	trades := models.TradeDao.FindTradesByHash(event.Hash)
	takerOrder := models.OrderDao.FindByID(trades[0].TakerOrderID)

	for _, trade := range trades {
		makerOrder := models.OrderDao.FindByID(trade.MakerOrderID)
		takerOrder.PendingAmount = takerOrder.PendingAmount.Sub(trade.Amount)
		makerOrder.PendingAmount = makerOrder.PendingAmount.Sub(trade.Amount)

		switch event.Status {
		case common.STATUS_FAILED:
			takerOrder.CanceledAmount = takerOrder.CanceledAmount.Add(trade.Amount)
			makerOrder.CanceledAmount = makerOrder.CanceledAmount.Add(trade.Amount)
		case common.STATUS_SUCCESSFUL:
			takerOrder.ConfirmedAmount = takerOrder.ConfirmedAmount.Add(trade.Amount)
			makerOrder.ConfirmedAmount = makerOrder.ConfirmedAmount.Add(trade.Amount)
		}

		makerOrder.AutoSetStatusByAmounts()
		_ = UpdateOrder(makerOrder)

		trade.Status = event.Status
		trade.ExecutedAt = time.Unix(int64(event.Timestamp), 0)
		_ = UpdateTrade(trade)
	}

	takerOrder.AutoSetStatusByAmounts()
	_ = UpdateOrder(takerOrder)

	return nil, nil
}

func NewMarketHandler(ctx context.Context, market *models.Market, engine *engine.Engine) (*MarketHandler, error) {
	orders := models.OrderDao.FindMarketPendingOrders(market.ID)

	// re-insert available orders into HydroEngine
	for _, order := range orders {
		if order.AvailableAmount.LessThanOrEqual(decimal.Zero) {
			continue
		}

		bookOrder := common.MemoryOrder{
			MarketID: order.MarketID,
			ID:       order.ID,
			Price:    order.Price,
			Amount:   order.AvailableAmount,
			Side:     order.Side,
		}
		msg := engine.ReInsertOrder(&bookOrder)
		pushMessage(msg)
	}

	marketHandler := MarketHandler{
		market:    market,
		eventChan: make(chan []byte),
		ctx:       ctx,

		hydroEngine: engine,
	}

	//// todo if Load Snapshot is necessary
	//res, err := kvStore.Get(common.GetMarketOrderbookSnapshotV2Key(market.ID))
	//if err == common.KVStoreEmpty {
	//	// do nothing
	//} else if err != nil {
	//	panic(fmt.Errorf("get snapshot error %v", err))
	//}
	//
	//var snapshot struct {
	//	Sequence uint64 `json:"sequence"`
	//}
	//
	//_ = json.Unmarshal([]byte(res), &snapshot)

	return &marketHandler, nil
}
