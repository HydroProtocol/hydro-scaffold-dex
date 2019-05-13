package api

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/shopspring/decimal"
)

type (
	BaseReq struct {
		Address string `json:"address"`
	}

	BaseResp struct {
		Status int    `json:"status"`
		Desc   string `json:"desc"`
	}

	MarketsReq struct {
		BaseReq
	}

	MarketsResp struct {
		BaseResp
		Markets []Market `json:"markets"`
		Count   int64    `json:"count"`
	}

	OrderBookReq struct {
		BaseReq
		MarketID string `json:"marketID" param:"marketID" validate:"required"`
	}

	OrderBookResp struct {
		BaseResp
		Data OrderBook
	}

	CandlesReq struct {
		BaseReq
		MarketID    string `json:"marketID"    param:"marketID"    validate:"required"`
		From        int64  `json:"from"        query:"from"        validate:"required"`
		To          int64  `json:"to"          query:"to"          validate:"required"`
		Granularity int64  `json:"granularity" query:"granularity" validate:"required"`
	}

	CandlesResp struct {
		BaseResp
		Data interface{}
	}

	QueryOrderReq struct {
		BaseReq
		MarketID string `json:"marketID" query:"marketID" validate:"required"`
		Status   string `json:"status"   query:"status"`
		Page     int    `json:"page"     query:"page"`
		PerPage  int    `json:"perPage"  query:"perPage"`
	}

	QueryOrderResp struct {
		Count  int64           `json:"count"`
		Orders []*models.Order `json:"orders"`
	}

	QuerySingleOrderReq struct {
		BaseReq
		OrderID string `json:"orderID" param:"orderID" validate:"required"`
	}

	QuerySingleOrderResp struct {
		Order *models.Order `json:"order"`
	}

	BuildOrderReq struct {
		BaseReq
		MarketID  string `json:"marketID"  validate:"required"`
		Side      string `json:"side"      validate:"required,oneof=buy sell"`
		OrderType string `json:"orderType" validate:"required,oneof=limit market"`
		Price     string `json:"price"     validate:"required"`
		Amount    string `json:"amount"    validate:"required"`
		Expires   int64  `json:"expires"`
	}

	BuildOrderResp struct {
		ID              string            `json:"id"`
		MarketID        string            `json:"marketID"`
		Side            string            `json:"side"`
		Type            string            `json:"type"`
		Price           decimal.Decimal   `json:"price"`
		Amount          decimal.Decimal   `json:"amount"`
		Json            *models.OrderJSON `json:"json"`
		AsMakerFeeRate  decimal.Decimal   `json:"asMakerFeeRate"`
		AsTakerFeeRate  decimal.Decimal   `json:"asTakerFeeRate"`
		MakerRebateRate decimal.Decimal   `json:"makerRebateRate"`
		GasFeeAmount    decimal.Decimal   `json:"gasFeeAmount"`
	}

	PlaceOrderReq struct {
		BaseReq
		ID        string `json:"orderID"   validate:"required,len=66"`
		Signature string `json:"signature" validate:"required"`
	}

	CancelOrderReq struct {
		BaseReq
		ID string `json:"id" param:"orderID" validate:"required,len=66"`
	}

	CacheOrder struct {
		OrderResponse         BuildOrderResp  `json:"orderResponse"`
		Address               string          `json:"address"`
		BalanceOfTokenToOffer decimal.Decimal `json:"balanceOfTokenToOffer"`
	}

	LockedBalanceReq struct {
		BaseReq
		//MarketID string `json:"marketID" validate:"required"`
	}

	LockedBalanceResp struct {
		LockedBalances []LockedBalance `json:"lockedBalances"`
	}

	QueryTradeReq struct {
		BaseReq
		MarketID string `json:"marketID" param:"marketID" validate:"required"`
		Status   string `json:"status"   query:"status"`
		Page     int    `json:"page"     query:"page"`
		PerPage  int    `json:"perPage"  query:"perPage"`
	}

	QueryTradeResp struct {
		Count  int64           `json:"count"`
		Trades []*models.Trade `json:"trades"`
	}

	FeesReq struct {
		BaseReq
		MarketID string `json:"marketID" query:"marketID" validate:"required"`
		Price    string `json:"price" query:"price"       validate:"required"`
		Amount   string `json:"amount" query:"amount"     validate:"required"`
	}

	FeesResp struct {
		GasFeeAmount          decimal.Decimal `json:"gasFeeAmount"`
		AsMakerTotalFeeAmount decimal.Decimal `json:"asMakerTotalFeeAmount"`
		AsMakerTradeFeeAmount decimal.Decimal `json:"asMakerTradeFeeAmount"`
		AsMakerFeeRate        decimal.Decimal `json:"asMakerFeeRate"`
		AsTakerTotalFeeAmount decimal.Decimal `json:"asTakerTotalFeeAmount"`
		AsTakerTradeFeeAmount decimal.Decimal `json:"asTakerTradeFeeAmount"`
		AsTakerFeeRate        decimal.Decimal `json:"asTakerFeeRate"`
	}
)

type (
	Market struct {
		ID                     string          `json:"id"`
		BaseToken              string          `json:"baseToken"`
		BaseTokenProjectUrl    string          `json:"baseTokenProjectUrl"`
		BaseTokenName          string          `json:"baseTokenName"`
		BaseTokenDecimals      int             `json:"baseTokenDecimals"`
		BaseTokenAddress       string          `json:"baseTokenAddress"`
		QuoteToken             string          `json:"quoteToken"`
		QuoteTokenDecimals     int             `json:"quoteTokenDecimals"`
		QuoteTokenAddress      string          `json:"quoteTokenAddress"`
		MinOrderSize           decimal.Decimal `json:"minOrderSize"`
		PricePrecision         int             `json:"pricePrecision"`
		PriceDecimals          int             `json:"priceDecimals"`
		AmountDecimals         int             `json:"amountDecimals"`
		AsMakerFeeRate         decimal.Decimal `json:"asMakerFeeRate"`
		AsTakerFeeRate         decimal.Decimal `json:"asTakerFeeRate"`
		GasFeeAmount           decimal.Decimal `json:"gasFeeAmount"`
		SupportedOrderTypes    []string        `json:"supportedOrderTypes"`
		MarketOrderMaxSlippage decimal.Decimal `json:"marketOrderMaxSlippage"`
		MarketStatus
	}

	MarketStatus struct {
		LastPriceIncrease   decimal.Decimal `json:"lastPriceIncrease"`
		LastPrice           decimal.Decimal `json:"lastPrice"`
		Price24h            decimal.Decimal `json:"price24h"`
		Amount24h           decimal.Decimal `json:"amount24h"`
		QuoteTokenVolume24h decimal.Decimal `json:"quoteTokenVolume24h"`
	}

	OrderBook struct {
		Bids [][2]string `json:"bids"`
		Asks [][2]string `json:"asks"`
	}

	LockedBalance struct {
		Symbol        string          `json:"symbol"`
		LockedBalance decimal.Decimal `json:"lockedBalance"`
	}
)

func (b *BaseReq) GetAddress() string {
	return b.Address
}

func (b *BaseReq) SetAddress(address string) {
	b.Address = address
}

type Param interface {
	GetAddress() string
	SetAddress(address string)
}
