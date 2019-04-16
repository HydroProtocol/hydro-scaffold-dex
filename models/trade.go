package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/shopspring/decimal"
	"time"
)

type ITradeDao interface {
	FindTradesByMarket(pair string, startTime time.Time, endTime time.Time) []*Trade
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
	ID              int64           `json:"id"               db:"id" primaryKey:"true" autoIncrement:"true"`
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

var TradeDao ITradeDao

func init() {
	TradeDao = &tradeDao{}
}

type tradeDao struct {
}

func (d *tradeDao) FindTradesByHash(hash string) []*Trade {
	trades := []*Trade{}
	findAllBy(
		&trades,
		&OpEq{
			"transaction_hash", hash,
		},
		map[string]OrderByDirection{"created_at": OrderByAsc},
		-1,
		-1,
	)

	return trades
}

func (d *tradeDao) FindTradesByMarket(marketID string, startTime time.Time, endTime time.Time) []*Trade {
	trades := []*Trade{}

	findAllBy(
		&trades,
		whereAnd(
			&OpEq{"market_id", marketID},
			&OpEq{"status", common.STATUS_SUCCESSFUL},
			&OpGt{"executed_at", startTime},
			&OpLt{"executed_at", endTime},
		),
		map[string]OrderByDirection{"executed_at": OrderByDesc},
		-1,
		-1,
	)

	return trades
}

func (d *tradeDao) FindAllTrades(marketID string) (int64, []*Trade) {
	trades := []*Trade{}
	conditions := whereAnd(
		&OpEq{"market_id", marketID},
		&OpEq{"status", common.STATUS_SUCCESSFUL},
	)
	findAllBy(
		&trades,
		conditions,
		map[string]OrderByDirection{"created_at": OrderByAsc},
		-1,
		-1,
	)

	count := findCountBy(&Trade{}, conditions)
	return int64(count), trades
}

func (d *tradeDao) FindAccountMarketTrades(account, marketID, status string, limit, offset int) (int64, []*Trade) {
	trades := []*Trade{}
	conditions := whereAnd(
		&OpEq{"market_id", marketID},
		whereOr(
			&OpEq{"taker", account},
			&OpEq{"maker", account},
		),
	)

	findAllBy(
		&trades,
		conditions,
		map[string]OrderByDirection{"created_at": OrderByAsc},
		limit,
		offset,
	)

	count := findCountBy(&Trade{}, conditions)
	return int64(count), trades
}

func (d *tradeDao) InsertTrade(trade *Trade) error {
	id, err := insert(trade)

	if err != nil {
		return err
	}

	trade.ID = id

	return nil
}

func (*tradeDao) UpdateTrade(trade *Trade) error {
	return update(trade, "Status", "TransactionID", "TransactionHash", "ExecutedAt")
}

func (*tradeDao) FindTradeByID(id int64) *Trade {
	var trade Trade

	findBy(&trade, &OpEq{"id", id}, nil)

	empty := Trade{}
	if trade == empty {
		return nil
	}

	return &trade
}

func (*tradeDao) FindTradeByTransactionID(transactionID int64) []*Trade {
	trades := []*Trade{}

	findAllBy(
		&trades,
		&OpEq{"transaction_id", transactionID},
		map[string]OrderByDirection{"created_at": OrderByAsc},
		-1,
		-1,
	)

	return trades
}

func (*tradeDao) Count() int {
	sql := "select count(*) from trades"
	var count int
	err := DB.QueryRowx(sql).Scan(&count)

	if err != nil {
		panic(err)
	}

	return count
}
