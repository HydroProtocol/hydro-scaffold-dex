package api

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMarkets(t *testing.T) {
	models.MockMarketDao()
	models.MockTradeDao()

	url := "/markets"
	resp := request(url, "GET", "", nil)

	markets := resp.Data.(map[string]interface{})["markets"]
	if len(markets.([]interface{})) == 0 {
		t.Errorf("no markets error")
	}

	marketsWethDai := markets.([]interface{})[0]
	assert.EqualValues(t, "WETH-DAI", marketsWethDai.(map[string]interface{})["id"])
	marketsHotDai := markets.([]interface{})[1]
	assert.EqualValues(t, "HOT-DAI", marketsHotDai.(map[string]interface{})["id"])
}

func TestGetOrderBookAPI(t *testing.T) {
	models.MockMarketDao()
	mockSnapshot()
	url := "/markets/HOT-DAI/orderbook?address=0x5409ed021d9299bf6814279a6a1411a7e866a631"
	resp := request(url, "GET", "", nil)

	assert.EqualValues(t, 0, resp.Status)

	orderBook := resp.Data.(map[string]interface{})["orderBook"]
	assert.EqualValues(t, 0, orderBook.(map[string]interface{})["sequence"])
}

func TestGetOrderBook(t *testing.T) {
	mockSnapshot()
	baseReq := BaseReq{Address: "0x5409ed021d9299bf6814279a6a1411a7e866a631"}
	req := &OrderBookReq{
		MarketID: "HOT-DAI",
		BaseReq:  baseReq,
	}
	snapshot, _ := GetOrderBook(req)

	t.Log(snapshot)
	snapshot2 := SnapshotV2{
		Bids: [][2]string{{"1.3", "3.4"}, {"1.2", "3.4"}},
		Asks: [][2]string{{"1.4", "3.4"}, {"1.5", "3.4"}},
	}
	t.Log(snapshot2)
	reflect.DeepEqual(snapshot, snapshot2)
}

type MCache struct {
	mock.Mock
}

func (m *MCache) Set(key string, value string, expire time.Duration) error {
	args := m.Called(key, value, expire)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (m *MCache) Get(key string) (string, error) {
	args := m.Called(key)
	if args.Get(1) == nil {
		return args.Get(0).(string), nil
	}
	return args.Get(0).(string), args.Get(1).(error)
}

func (m *MCache) Push(key []byte) error {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (m *MCache) Pop() ([]byte, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return args.Get(0).([]byte), nil
	}
	return args.Get(0).([]byte), args.Get(1).(error)
}

func mockSnapshot() {
	book := common.NewOrderbook("HOT-DAI")
	book.InsertOrder(NewLimitOrder("o1", "buy", "1.2", "3.4"))
	book.InsertOrder(NewLimitOrder("o2", "buy", "1.3", "3.4"))
	book.InsertOrder(NewLimitOrder("o3", "sell", "1.4", "3.4"))
	book.InsertOrder(NewLimitOrder("o4", "sell", "1.5", "3.4"))

	snapshot := book.SnapshotV2()

	bts, err := json.Marshal(snapshot)

	if err != nil {
		panic(err)
	}

	cacheService := &MCache{}
	CacheService = cacheService

	cacheService.On("Set", mock.Anything, mock.Anything, mock.Anything).Times(10).Return(nil)
	cacheService.On("Get", mock.Anything).Return(string(bts), nil).Once()
	cacheService.On("Push", mock.Anything).Times(10).Return(nil)
	cacheService.On("Pop").Times(1000).Return([]string{}, nil)
}

// NewLimitOrder ...
func NewLimitOrder(id string, side string, price string, amount string) *common.MemoryOrder {
	return NewOrder(id, side, price, amount, "limit")
}

func NewOrder(id, side, price, amount, _type string) *common.MemoryOrder {
	if len(id) <= 0 {
		panic(fmt.Errorf("ID can't be blank"))
	}

	amountDecimal, err := decimal.NewFromString(amount)

	if side != "buy" && side != "sell" {
		panic(fmt.Errorf("side should be buy/sell. passed: %s", side))
	}

	if err != nil {
		panic(fmt.Errorf("amount decimal error, Amount: %s, error: %+v", amount, err))
	}

	priceDecimal, err := decimal.NewFromString(price)
	if err != nil {
		panic(fmt.Errorf("price decimal error, Price: %s, error: %+v", price, err))
	}

	return &common.MemoryOrder{
		ID:    id,
		Side:  side,
		Price: priceDecimal,

		Amount: amountDecimal,
		Type:   _type,
	}
}
