package models

import (
	"github.com/shopspring/decimal"
)

type IMarketDao interface {
	FindAllMarkets() []*Market
	FindPublishedMarkets() []*Market
	FindMarketByID(marketID string) *Market
	InsertMarket(market *Market) error
	UpdateMarket(market *Market) error
}

type Market struct {
	ID                string `json:"id"                db:"id" primaryKey:"true" gorm:"primary_key"`
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
	IsPublished       bool            `json:"isPublished"       db:"is_published"`
}

func (Market) TableName() string {
	return "markets"
}

var MarketDao IMarketDao
var MarketDaoPG IMarketDao

func init() {
	MarketDao = &marketDaoPG{}
	MarketDaoPG = MarketDao
}

type marketDaoPG struct {
}

func (marketDaoPG) FindPublishedMarkets() []*Market {
	var markets []*Market
	DB.Where("is_published = ?", true).Find(&markets)
	return markets
}

func (marketDaoPG) FindAllMarkets() []*Market {
	var markets []*Market
	DB.Find(&markets)
	return markets
}

func (marketDaoPG) FindMarketByID(marketID string) *Market {
	var market Market
	DB.Where("id = ?", marketID).First(&market)
	if market.ID == "" {
		return nil
	}
	return &market
}

func (marketDaoPG) InsertMarket(market *Market) error {
	return DB.Create(market).Error
}

func (marketDaoPG) UpdateMarket(market *Market) error {
	return DB.Save(market).Error
}
