package models

import (
	"github.com/shopspring/decimal"
)

type IMarketDao interface {
	FindAllMarkets() []*Market
	FindMarketByID(marketID string) *Market
	InsertMarket(market *Market) error
	UpdateMarket(market *Market) error
}

type Market struct {
	ID                string `json:"id"                db:"id"`
	BaseTokenSymbol   string `json:"baseTokenSymbol"   db:"base_token_symbol"`
	BaseTokenName     string `json:"BaseTokenName"     db:"base_token_name"`
	BaseTokenAddress  string `json:"baseTokenAddress"  db:"base_token_address"`
	BaseTokenDecimals int    `json:"baseTokenDecimals" db:"base_token_decimals"`

	QuoteTokenSymbol   string `json:"quoteTokenSymbol"   db:"quote_token_symbol"`
	QuoteTokenName     string `json:"QuoteTokenName"     db:"quote_token_name"`
	QuoteTokenAddress  string `json:"quoteTokenAddress"  db:"quote_token_address"`
	QuoteTokenDecimals int    `json:"quoteTokenDecimals" db:"quote_token_decimals"`

	MinOrderSize      decimal.Decimal `json:"minOrderSize"      db:"min_order_size"`
	PricePrecision    int             `json:"pricePrecision"    db:"price_precision"`
	PriceDecimals     int             `json:"priceDecimals"     db:"price_decimals"`
	AmountDecimals    int             `json:"amountDecimals"    db:"amount_decimals"`
	MakerFeeRate      decimal.Decimal `json:"makerFeeRate"      db:"maker_fee_rate"`
	TakerFeeRate      decimal.Decimal `json:"takerFeeRate"      db:"taker_fee_rate"`
	GasUsedEstimation int             `json:"gasUsedEstimation" db:"gas_used_estimation"`
}

var MarketDao IMarketDao

func init() {
	MarketDao = &marketDao{}
}

type marketDao struct {
}

func (marketDao) FindAllMarkets() []*Market {
	markets := []*Market{}
	findAllBy(&markets, nil, nil, -1, -1)
	return markets
}

func (marketDao) FindMarketByID(marketID string) *Market {
	var market Market
	findBy(&market, &OpEq{"id", marketID}, nil)
	if market.ID == "" {
		return nil
	}
	return &market
}

func (marketDao) InsertMarket(market *Market) error {
	_, err := insert(market)
	return err
}

func (marketDao) UpdateMarket(market *Market) error {
	return update(market, "MinOrderSize", "PricePrecision", "PriceDecimals", "AmountDecimals", "MakerFeeRate", "TakerFeeRate", "GasUsedEstimation")
}
