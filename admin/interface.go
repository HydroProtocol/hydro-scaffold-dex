package main

import (
	"encoding/json"
	"fmt"
	"github.com/HydroProtocol/hydro-box-dex/models"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"net/url"
)

type IAdminApi interface {
	NewMarket(market string) error
	EditMarket(market string) error
	ListAccountOrders(account string) error
	CancelOrder(ID string) error
}

type Admin struct {
	client         utils.IHttpClient
	AdminApiUrl    string
	MarketUrl      string
	CancelOrderUrl string
	ListOrderUrl   string
}

func NewAdmin(adminApiUrl string) IAdminApi {
	_, err := url.Parse(adminApiUrl)
	if err != nil {
		panic(err)
	}

	a := Admin{}
	a.client = utils.NewHttpClient(nil)
	a.AdminApiUrl = adminApiUrl
	a.MarketUrl = fmt.Sprintf("%s/%s", adminApiUrl, "markets")
	a.CancelOrderUrl = fmt.Sprintf("%s/%s", adminApiUrl, "orders")
	a.ListOrderUrl = fmt.Sprintf("%s/%s", adminApiUrl, "orders")

	return &a
}

func (a *Admin) NewMarket(marketData string) (err error) {
	var market models.Market
	err = json.Unmarshal([]byte(marketData), &market)
	if err != nil {
		return
	}

	err, _, _ = a.client.Post(a.MarketUrl, nil, market, nil)
	return
}

func (a *Admin) EditMarket(marketData string) (err error) {
	var market models.Market
	err = json.Unmarshal([]byte(marketData), &market)
	if err != nil {
		return
	}

	err, _, _ = a.client.Put(a.MarketUrl, nil, market, nil)
	return
}

func (a *Admin) ListAccountOrders(reqData string) (err error) {
	var req = struct {
		Trader   string `json:"trader"`
		MarketID string `json:"market_id"`
		Status   string `json:"status"`
		Offset   int    `json:"offset"`
		Limit    int    `json:"limit"`
	}{}

	err = json.Unmarshal([]byte(reqData), &req)
	if err != nil {
		return
	}
	var params []utils.KeyValue
	params = append(params, utils.KeyValue{Key: "account", Value: req.Trader})
	params = append(params, utils.KeyValue{Key: "marketID", Value: req.MarketID})
	params = append(params, utils.KeyValue{Key: "status", Value: req.Status})
	params = append(params, utils.KeyValue{Key: "offset", Value: fmt.Sprintf("%d", req.Offset)})
	params = append(params, utils.KeyValue{Key: "limit", Value: fmt.Sprintf("%d", req.Limit)})

	err, _, _ = a.client.Get(a.ListOrderUrl, params, nil, nil)
	return
}

func (a *Admin) CancelOrder(ID string) (err error) {
	err, _, _ = a.client.Delete(fmt.Sprintf("%s/%s", a.CancelOrderUrl, ID), nil, nil, nil)
	return
}
