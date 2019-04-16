package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
	"math/rand"
	"time"
)

func GetLockedBalance(p Param) (interface{}, error) {
	req := p.(*LockedBalanceReq)
	tokens := models.TokenDao.GetAllTokens()

	var lockedBalances []LockedBalance

	for _, token := range tokens {
		lockedBalance := models.BalanceDao.GetByAccountAndSymbol(req.Address, token.Symbol, token.Decimals)
		lockedBalances = append(lockedBalances, LockedBalance{
			Symbol:        token.Symbol,
			LockedBalance: lockedBalance,
		})
	}

	return &LockedBalanceResp{
		LockedBalances: lockedBalances,
	}, nil
}

func GetOrder(p Param) (interface{}, error) {
	req := p.(*QueryOrderReq)
	if req.Status == "" {
		req.Status = common.ORDER_PENDING
	}
	if req.PerPage <= 0 {
		req.PerPage = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}

	offset := req.PerPage * (req.Page - 1)
	limit := req.PerPage

	count, orders := models.OrderDao.FindByAccount(req.Address, req.MarketID, req.Status, limit, offset)

	return &QueryOrderResp{
		Count:  count,
		Orders: orders,
	}, nil
}

func CancelOrder(p Param) (interface{}, error) {
	req := p.(*CancelOrderReq)
	order := models.OrderDao.FindByID(req.ID)
	if order == nil {
		return nil, NewApiError(-1, fmt.Sprintf("order %s not exist", req.ID))
	}

	cancelOrderEvent := common.CancelOrderEvent{
		Event: common.Event{
			Type:     common.EventCancelOrder,
			MarketID: order.MarketID,
		},
		Price: order.Price.String(),
		Side:  order.Side,
		ID:    order.ID,
	}

	return nil, QueueService.Push([]byte(utils.ToJsonString(cancelOrderEvent)))
}

func BuildOrder(p Param) (interface{}, error) {
	req := p.(*BuildOrderReq)
	err := checkBalanceAndAllowance(req, req.Address)
	if err != nil {
		return nil, err
	}

	buildOrderResponse, err := BuildAndCacheOrder(req.Address, req)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"order": buildOrderResponse,
	}, nil
}

func PlaceOrder(p Param) (interface{}, error) {
	order := p.(*PlaceOrderReq)
	if valid := hydro.IsValidOrderSignature(order.Address, order.ID, order.Signature); !valid {
		utils.Info("valid is %v", valid)
		return nil, errors.New("bad signature")
	}

	cacheOrder := getCacheOrderByOrderID(order.ID)

	if cacheOrder == nil {
		return nil, errors.New("place order error, please retry later")
	}

	cacheOrder.OrderResponse.Json.Signature = order.Signature

	ret := models.Order{
		ID:              order.ID,
		TraderAddress:   order.Address,
		MarketID:        cacheOrder.OrderResponse.MarketID,
		Side:            cacheOrder.OrderResponse.Side,
		Price:           cacheOrder.OrderResponse.Price,
		Amount:          cacheOrder.OrderResponse.Amount,
		Status:          common.ORDER_PENDING,
		Type:            cacheOrder.OrderResponse.Type,
		Version:         "hydro-v1",
		AvailableAmount: cacheOrder.OrderResponse.Amount,
		ConfirmedAmount: decimal.Zero,
		CanceledAmount:  decimal.Zero,
		PendingAmount:   decimal.Zero,
		MakerFeeRate:    cacheOrder.OrderResponse.AsMakerFeeRate,
		TakerFeeRate:    cacheOrder.OrderResponse.AsTakerFeeRate,
		MakerRebateRate: cacheOrder.OrderResponse.MakerRebateRate,
		GasFeeAmount:    cacheOrder.OrderResponse.GasFeeAmount,
		JSON:            utils.ToJsonString(cacheOrder.OrderResponse.Json),
		CreatedAt:       time.Now(),
	}

	newOrderEvent, _ := json.Marshal(common.NewOrderEvent{
		Event: common.Event{
			MarketID: cacheOrder.OrderResponse.MarketID,
			Type:     common.EventNewOrder,
		},
		Order: utils.ToJsonString(ret),
	})

	err := QueueService.Push(newOrderEvent)

	if err != nil {
		return nil, errors.New("place order failed, place try again")
	} else {
		return nil, nil
	}
}

func getCacheOrderByOrderID(orderID string) *CacheOrder {
	cacheOrderStr, err := CacheService.Get(generateOrderCacheKey(orderID))

	if err != nil {
		utils.Error("get cache order error: %v", err)
		return nil
	}

	var cacheOrder CacheOrder

	err = json.Unmarshal([]byte(cacheOrderStr), &cacheOrder)
	if err != nil {
		utils.Error("get cache order error: %v, cache order is: %v", err, cacheOrderStr)
		return nil
	}

	return &cacheOrder
}

func checkBalanceAndAllowance(order *BuildOrderReq, address string) error {
	market := models.MarketDao.FindMarketByID(order.MarketID)
	if market == nil {
		return MarketNotFoundError(order.MarketID)
	}

	minPriceUnit := decimal.New(1, int32(-1*market.PriceDecimals))

	price := utils.StringToDecimal(order.Price)

	if price.LessThanOrEqual(decimal.Zero) {
		return NewApiError(-1, "invalid_price")
	}

	if !price.Mod(minPriceUnit).Equal(decimal.Zero) {
		return NewApiError(-1, "invalid_price_unit")
	}

	minAmountUnit := decimal.New(1, int32(-1*market.AmountDecimals))

	amount := utils.StringToDecimal(order.Amount)

	if amount.LessThanOrEqual(decimal.Zero) {
		return NewApiError(-1, "invalid_amount")
	}

	if !amount.Mod(minAmountUnit).Equal(decimal.Zero) {
		return NewApiError(-1, "invalid_amount_unit")
	}

	baseTokenLockedBalance := models.BalanceDao.GetByAccountAndSymbol(address, market.BaseTokenSymbol, market.BaseTokenDecimals)
	baseTokenBalance := hydro.GetTokenBalance(market.BaseTokenAddress, address)
	baseTokenAllowance := hydro.GetTokenAllowance(market.BaseTokenAddress, config.Getenv("HSK_PROXY_ADDRESS"), address)

	quoteTokenLockedBalance := models.BalanceDao.GetByAccountAndSymbol(address, market.QuoteTokenSymbol, market.QuoteTokenDecimals)
	quoteTokenBalance := hydro.GetTokenBalance(market.QuoteTokenAddress, address)
	quoteTokenAllowance := hydro.GetTokenAllowance(market.QuoteTokenAddress, config.Getenv("HSK_PROXY_ADDRESS"), address)

	var quoteTokenHugeAmount decimal.Decimal
	var baseTokenHugeAmount decimal.Decimal

	feeDetail := calculateFee(price, amount, market, address)
	feeAmount := feeDetail.AsTakerTotalFeeAmount

	quoteTokenHugeAmount = amount.Mul(decimal.New(1, int32(market.QuoteTokenDecimals))).Mul(price)
	baseTokenHugeAmount = amount.Mul(decimal.New(1, int32(market.QuoteTokenDecimals)))

	if order.Side == "sell" {
		if quoteTokenHugeAmount.LessThanOrEqual(feeAmount) {
			return NewApiError(-1, fmt.Sprintf("amount: %s less than fee: %s", quoteTokenHugeAmount.String(), feeAmount.String()))
		}

		availableBaseTokenAmount := baseTokenBalance.Sub(baseTokenLockedBalance)
		if baseTokenHugeAmount.GreaterThan(availableBaseTokenAmount) {
			return NewApiError(-1, fmt.Sprintf("%s balance not enough, available balance is %s, require amount is %s", market.BaseTokenSymbol, availableBaseTokenAmount.String(), baseTokenHugeAmount.String()))
		}

		if baseTokenHugeAmount.GreaterThan(baseTokenAllowance) {
			return NewApiError(-1, fmt.Sprintf("%s allowance not enough, allowance is %s, require amount is %s", market.BaseTokenSymbol, baseTokenAllowance.String(), baseTokenHugeAmount.String()))
		}
	} else {
		availableQuoteTokenAmount := quoteTokenBalance.Sub(quoteTokenLockedBalance)
		requireAmount := quoteTokenHugeAmount.Add(feeAmount)
		if requireAmount.GreaterThan(availableQuoteTokenAmount) {
			return NewApiError(-1, fmt.Sprintf("%s balance not enough, available balance is %s, require amount is %s", market.QuoteTokenSymbol, availableQuoteTokenAmount.String(), requireAmount.String()))
		}

		if requireAmount.GreaterThan(quoteTokenAllowance) {
			return NewApiError(-1, fmt.Sprintf("%s allowance not enough, available balance is %s, require amount is %s", market.QuoteTokenSymbol, quoteTokenAllowance.String(), requireAmount.String()))
		}
	}

	// will add check of precision later

	return nil
}

func BuildAndCacheOrder(address string, order *BuildOrderReq) (*BuildOrderResp, error) {
	market := models.MarketDao.FindMarketByID(order.MarketID)
	amount := utils.StringToDecimal(order.Amount)
	price := utils.StringToDecimal(order.Price)

	fee := calculateFee(price, amount, market, address)

	gasFeeInQuoteToken := fee.GasFeeAmount
	makerRebateRate := decimal.Zero
	offeredAmount := decimal.Zero

	var baseTokenHugeAmount decimal.Decimal
	var quoteTokenHugeAmount decimal.Decimal

	baseTokenHugeAmount = amount.Mul(decimal.New(1, int32(market.BaseTokenDecimals)))
	quoteTokenHugeAmount = price.Mul(amount).Mul(decimal.New(1, int32(market.QuoteTokenDecimals)))

	orderData := hydro.GenerateOrderData(
		int64(2),
		getExpiredAt(order.Expires),
		rand.Int63(),
		market.TakerFeeRate,
		market.MakerFeeRate,
		decimal.Zero,
		order.Side == "sell",
		order.OrderType == "market",
		false)

	orderJson := models.OrderJSON{
		Trader:                  address,
		Relayer:                 config.Getenv("HSK_RELAYER_ADDRESS"),
		BaseCurrency:            market.BaseTokenAddress,
		QuoteCurrency:           market.QuoteTokenAddress,
		BaseCurrencyHugeAmount:  baseTokenHugeAmount,
		QuoteCurrencyHugeAmount: quoteTokenHugeAmount,
		GasTokenHugeAmount:      gasFeeInQuoteToken,
		Data:                    orderData,
	}

	sdkOrder := sdk.NewOrderWithData(address,
		config.Getenv("HSK_RELAYER_ADDRESS"),
		market.BaseTokenAddress,
		market.QuoteTokenAddress,
		utils.DecimalToBigInt(baseTokenHugeAmount),
		utils.DecimalToBigInt(quoteTokenHugeAmount),
		utils.DecimalToBigInt(gasFeeInQuoteToken),
		orderData,
		"",
	)

	orderHash := hydro.GetOrderHash(sdkOrder)
	orderResponse := BuildOrderResp{
		ID:              utils.Bytes2HexP(orderHash),
		Json:            &orderJson,
		Side:            order.Side,
		Type:            order.OrderType,
		Price:           price,
		Amount:          amount,
		MarketID:        order.MarketID,
		AsMakerFeeRate:  market.MakerFeeRate,
		AsTakerFeeRate:  market.TakerFeeRate,
		MakerRebateRate: makerRebateRate,
		GasFeeAmount:    gasFeeInQuoteToken,
	}

	cacheOrder := CacheOrder{
		OrderResponse:         orderResponse,
		Address:               address,
		BalanceOfTokenToOffer: offeredAmount,
	}

	// Cache the build order for 60 seconds, if we still not get signature in the period. The order will be dropped.
	err := CacheService.Set(generateOrderCacheKey(orderResponse.ID), utils.ToJsonString(cacheOrder), time.Second*60)
	return &orderResponse, err
}

func generateOrderCacheKey(orderID string) string {
	return "OrderCache:" + orderID
}

func getExpiredAt(expiresInSeconds int64) int64 {
	if time.Duration(expiresInSeconds)*time.Second > time.Hour {
		return time.Now().Unix() + expiresInSeconds
	} else {
		return time.Now().Unix() + 60*60*24*365*100
	}
}

func isMarketBuyOrder(order *BuildOrderReq) bool {
	return order.OrderType == "market" && order.Side == "buy"
}

func isMarketOrder(order *BuildOrderReq) bool {
	return order.OrderType == "market"
}
