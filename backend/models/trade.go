package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/shopspring/decimal"
	"time"
)

type ITradeDao interface {
	FindTradesByMarket(marketID string, startTime time.Time, endTime time.Time) []*Trade
	FindAllTrades(marketID string) (int64, []*Trade)
	FindTradesByHash(hash string) []*Trade
	FindTradeByID(id int64) *Trade
	FindAccountMarketTrades(account, marketID, status string, limit, offset int) (int64, []*Trade)

	InsertTrade(trade *Trade) error
	UpdateTrade(trade *Trade) error
	Count() int
	FindTradeByTransactionID(transactionID int64) []*Trade
}

type Trade struct {
	ID              int64           `json:"id"               db:"id" primaryKey:"true" autoIncrement:"true" gorm:"primary_key"`
	TransactionID   int64           `json:"transactionID"    db:"transaction_id"`
	TransactionHash string          `json:"transactionHash"  db:"transaction_hash"`
	Status          string          `json:"status"           db:"status"`
	MarketID        string          `json:"marketID"         db:"market_id"`
	Maker           string          `json:"maker"            db:"maker"`
	Taker           string          `json:"taker"            db:"taker"`
	TakerSide       string          `json:"takerSide"        db:"taker_side"`
	MakerOrderID    string          `json:"makerOrderID"     db:"maker_order_id"`
	TakerOrderID    string          `json:"takerOrderID"     db:"taker_order_id"`
	Sequence        int             `json:"sequence"         db:"sequence"`
	Amount          decimal.Decimal `json:"amount"           db:"amount"`
	Price           decimal.Decimal `json:"price"            db:"price"`
	ExecutedAt      time.Time       `json:"executedAt"       db:"executed_at"`
	CreatedAt       time.Time       `json:"createdAt"        db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt"        db:"updated_at"`
}

func (Trade) TableName() string {
	return "trades"
}

var TradeDao ITradeDao
var TradeDaoPG ITradeDao

func init() {
	TradeDao = &tradeDaoPG{}
	TradeDaoPG = TradeDao
}

type tradeDaoPG struct {
}

func (tradeDaoPG) FindTradesByMarket(marketID string, startTime time.Time, endTime time.Time) []*Trade {
	var trades []*Trade

	DB.Where("market_id = ? and status = ? and executed_at between ? and ? ", marketID, common.STATUS_SUCCESSFUL, startTime, endTime).Order("executed_at desc").Find(&trades)
	return trades
}

func (tradeDaoPG) FindAllTrades(marketID string) (int64, []*Trade) {
	var trades []*Trade
	var count int64

	DB.Where("market_id = ? and status = ?", marketID, common.STATUS_SUCCESSFUL).Order("created_at desc").Find(&trades).Count(&count)
	return count, trades
}

func (tradeDaoPG) FindTradesByHash(hash string) []*Trade {
	var trades []*Trade
	DB.Where("transaction_hash = ?", hash).Order("created_at desc").Find(&trades)
	return trades
}

func (tradeDaoPG) FindTradeByID(id int64) *Trade {
	var trade Trade

	DB.Where("id = ?", id).Find(&trade)
	if trade.Status == "" {
		return nil
	}

	return &trade
}

func (tradeDaoPG) FindAccountMarketTrades(account, marketID, status string, limit, offset int) (int64, []*Trade) {
	var trades []*Trade
	var count int64

	DB.Where("market_id = ? and (taker = ? or maker = ?)", marketID, account, account).Order("created_at desc").Find(&trades).Count(&count)
	return count, trades
}

func (tradeDaoPG) InsertTrade(trade *Trade) error {
	return DB.Create(trade).Error
}

func (tradeDaoPG) UpdateTrade(trade *Trade) error {
	return DB.Save(trade).Error
}

func (tradeDaoPG) Count() int {
	var count int
	DB.Model(&Trade{}).Count(&count)
	return count
}

func (tradeDaoPG) FindTradeByTransactionID(transactionID int64) []*Trade {
	var trades []*Trade

	DB.Where("transaction_id = ? ", transactionID).Order("created_at asc").Find(&trades)
	return trades
}
