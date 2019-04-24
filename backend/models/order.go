package models

import (
	"encoding/json"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
	"time"
)

type IOrderDao interface {
	FindMarketPendingOrders(marketID string) []*Order
	FindByAccount(trader, marketID, status string, offset, limit int) (int64, []*Order)
	FindByID(id string) *Order
	InsertOrder(order *Order) error
	UpdateOrder(order *Order) error
	Count() int
}

var OrderDaoSqlite IOrderDao
var OrderDaoPG IOrderDao

type orderDao struct {
}

func init() {
	OrderDaoSqlite = &orderDao{}
	OrderDaoPG = &orderDaoPG{}
}

type Order struct {
	ID              string          `json:"id" db:"id" primaryKey:"true" gorm:"primary_key"`
	TraderAddress   string          `json:"traderAddress" db:"trader_address"`
	MarketID        string          `json:"marketID" db:"market_id"`
	Side            string          `json:"side" db:"side"`
	Price           decimal.Decimal `json:"price" db:"price"`
	Amount          decimal.Decimal `json:"amount" db:"amount"`
	Status          string          `json:"status" db:"status"`
	Type            string          `json:"type" db:"type"`
	Version         string          `json:"version" db:"version"`
	AvailableAmount decimal.Decimal `json:"availableAmount" db:"available_amount"`
	ConfirmedAmount decimal.Decimal `json:"confirmedAmount" db:"confirmed_amount"`
	CanceledAmount  decimal.Decimal `json:"canceledAmount" db:"canceled_amount"`
	PendingAmount   decimal.Decimal `json:"pendingAmount" db:"pending_amount"`
	MakerFeeRate    decimal.Decimal `json:"makerFeeRate" db:"maker_fee_rate"`
	TakerFeeRate    decimal.Decimal `json:"takerFeeRate" db:"taker_fee_rate"`
	MakerRebateRate decimal.Decimal `json:"makerRebateRate" db:"maker_rebate_rate"`
	GasFeeAmount    decimal.Decimal `json:"gasFeeAmount" db:"gas_fee_amount"`
	JSON            string          `json:"json" db:"json"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
}

func (o *Order) AutoSetStatusByAmounts() {
	if o.ConfirmedAmount.Equal(o.Amount) {
		o.Status = common.ORDER_FULL_FILLED
	} else if o.CanceledAmount.Equal(o.Amount) {
		o.Status = common.ORDER_CANCELED
	} else if o.AvailableAmount.Add(o.PendingAmount).GreaterThan(decimal.Zero) {
		o.Status = common.ORDER_PENDING
	} else {
		o.Status = common.ORDER_PARTIAL_FILLED
	}
}

type OrderJSON struct {
	Trader                  string          `json:"trader"`
	Relayer                 string          `json:"relayer"`
	BaseCurrencyHugeAmount  decimal.Decimal `json:"baseTokenAmount"`
	QuoteCurrencyHugeAmount decimal.Decimal `json:"quoteTokenAmount"`
	BaseCurrency            string          `json:"baseToken"`
	QuoteCurrency           string          `json:"quoteToken"`
	GasTokenHugeAmount      decimal.Decimal `json:"gasTokenAmount"`
	Signature               string          `json:"signature"`
	Data                    string          `json:"data"`
}

type ECSignature struct {
	Config string `json:"config"`
	R      string `json:"r"`
	S      string `json:"s"`
}

func (o Order) GetOrderJson() *OrderJSON {
	var orderJson OrderJSON
	json.Unmarshal([]byte(o.JSON), &orderJson)
	return &orderJson
}

func (*orderDao) Count() int {
	sql := "select count(*) from orders"
	var count int
	err := DBSqlite.QueryRowx(sql).Scan(&count)

	if err != nil {
		panic(err)
	}

	return count
}

func (*orderDao) FindMarketPendingOrders(marketID string) []*Order {
	orders := []*Order{}

	findAllBy(
		&orders,
		whereAnd(
			&OpEq{"status", common.ORDER_PENDING},
			&OpEq{"market_id", marketID},
		),
		map[string]OrderByDirection{"created_at": OrderByAsc},
		-1,
		-1,
	)

	return orders
}

func (*orderDao) FindByAccount(trader, marketID, status string, limit, offset int) (int64, []*Order) {
	orders := []*Order{}

	conditions := whereAnd(
		&OpEq{"trader_address", trader},
		&OpEq{"market_id", marketID},
		&OpEq{"status", status},
	)

	findAllBy(
		&orders,
		conditions,
		map[string]OrderByDirection{"created_at": OrderByAsc},
		limit,
		offset,
	)

	count := findCountBy(&Order{}, conditions)

	return int64(count), orders
}

func (*orderDao) FindByID(id string) *Order {
	var order Order

	findBy(&order, &OpEq{"id", id}, nil)

	if order.ID == "" {
		return nil
	}

	return &order
}

func (*orderDao) InsertOrder(order *Order) error {
	_, err := insert(order)
	return err
}

func (*orderDao) UpdateOrder(order *Order) error {
	return update(order, "AvailableAmount", "ConfirmedAmount", "CanceledAmount", "PendingAmount", "Status")
}

//pg
type orderDaoPG struct {
}

func (Order) TableName() string {
	return "orders"
}

func (orderDaoPG) FindMarketPendingOrders(marketID string) (orders []*Order) {
	DBPG.Where("status = 'pending' and market_id = ?", marketID).Order("created_at asc").Find(&orders)
	return
}

func (orderDaoPG) FindByAccount(trader, marketID, status string, offset, limit int) (count int64, orders []*Order) {
	DBPG.Where("trader_address = ? and market_id = ? and status = ?", trader, marketID, status).Order("created_at desc").Limit(limit).Offset(offset).Find(&orders)
	DBPG.Model(&Order{}).Where("trader_address = ? and market_id = ? and status = ?", trader, marketID, status).Count(&count)
	return
}

func (orderDaoPG) FindByID(id string) *Order {
	var order Order
	DBPG.Where("id = ?", id).First(&order)
	if order.ID == "" {
		return nil
	}
	return &order
}

func (orderDaoPG) InsertOrder(order *Order) error {
	return DBPG.Create(order).Error
}

func (orderDaoPG) UpdateOrder(order *Order) error {
	return DBPG.Save(order).Error
}

func (o orderDaoPG) Count() (count int) {
	err := DBPG.Model(&Order{}).Count(&count).Error
	if err != nil {
		utils.Error("count orders error: %v", err)
	}

	return
}
