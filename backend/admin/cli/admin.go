package admincli

import (
	"encoding/json"
	"fmt"
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"net/url"
)

const (
	DefaultMinOrderSize      = "0.01"
	DefaultPricePrecision    = 5
	DefaultPriceDecimals     = 5
	DefaultAmountDecimals    = 5
	DefaultMakerFeeRate      = "0.01"
	DefaultTakerFeeRate      = "0.03"
	DefaultGasUsedEstimation = 190000

	DefaultLimit  = "10"
	DefaultOffset = "10"
	DefaultStatus = "pending"

	DefaultAdminAPIURL = "http://localhost:3003"
)

type IAdminApi interface {
	Status() error

	NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation string) error
	ListMarkets() error
	UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish string) error
	PublishMarket(marketID string) error
	UnPublishMarket(marketID string) error
	UpdateMarketFee(marketID, makerFee, takerFee string) error

	ListAccountOrders(address, limit, offset, status string) error
	ListAccountBalances(address, limit, offset string) error
	ListAccountTrades(address, limit, offset, status string) error

	CancelOrder(ID string) error

	RestartEngine() error
}

type Admin struct {
	client           utils.IHttpClient
	erc20            ethereum.IErc20
	AdminApiUrl      string
	MarketUrl        string
	CancelOrderUrl   string
	ListOrderUrl     string
	ListBalanceUrl   string
	ListTradeUrl     string
	RestartEngineUrl string
	StatusUrl        string
}

func NewAdmin(adminApiUrl string, httpClient utils.IHttpClient, erc20 ethereum.IErc20) IAdminApi {
	if len(adminApiUrl) == 0 {
		adminApiUrl = DefaultAdminAPIURL
	} else {
		_, err := url.Parse(adminApiUrl)
		if err != nil {
			panic(err)
		}
	}

	a := Admin{}
	if httpClient == nil {
		a.client = utils.NewHttpClient(nil)
	} else {
		a.client = httpClient
	}

	if erc20 == nil {
		a.erc20 = ethereum.NewErc20Service(nil)
	} else {
		a.erc20 = erc20
	}

	a.AdminApiUrl = adminApiUrl
	a.MarketUrl = fmt.Sprintf("%s/%s", adminApiUrl, "markets")
	a.CancelOrderUrl = fmt.Sprintf("%s/%s", adminApiUrl, "orders")
	a.ListOrderUrl = fmt.Sprintf("%s/%s", adminApiUrl, "orders")
	a.ListTradeUrl = fmt.Sprintf("%s/%s", adminApiUrl, "trades")
	a.ListBalanceUrl = fmt.Sprintf("%s/%s", adminApiUrl, "balances")
	a.RestartEngineUrl = fmt.Sprintf("%s/%s", adminApiUrl, "restart_engine")
	a.StatusUrl = fmt.Sprintf("%s/%s", adminApiUrl, "status")

	return &a
}

func (a *Admin) Status() (err error) {
	var statusResult struct {
		Data map[string]string `json:"data"`
	}

	err, _, data := a.client.Get(a.StatusUrl, nil, nil, nil)

	if err != nil {
		utils.Debug("call status error: %v", err)
		return
	}

	err = json.Unmarshal(data, &statusResult)
	return
}

func (a *Admin) ListMarkets() (err error) {
	err, _, _ = a.client.Get(a.MarketUrl, nil, nil, nil)
	return
}

func (a *Admin) NewMarket(marketID, baseTokenAddress, quoteTokenAddress, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation string) (err error) {
	err, baseTokenSymbol := a.erc20.Name(baseTokenAddress)
	if err != nil || len(baseTokenSymbol) == 0 {
		return
	}

	err, baseTokenName := a.erc20.Name(baseTokenAddress)
	if err != nil || len(baseTokenName) == 0 {
		return
	}

	err, baseTokenDecimals := a.erc20.Decimals(baseTokenAddress)
	if err != nil {
		return
	}

	err, quoteTokenSymbol := a.erc20.Symbol(quoteTokenAddress)
	if err != nil || len(quoteTokenSymbol) == 0 {
		return
	}

	err, quoteTokenName := a.erc20.Name(quoteTokenAddress)
	if err != nil || len(quoteTokenName) == 0 {
		return
	}

	err, quoteTokenDecimals := a.erc20.Decimals(quoteTokenAddress)
	if err != nil {
		return
	}

	market := models.Market{
		ID: marketID,

		BaseTokenAddress:  baseTokenAddress,
		BaseTokenName:     baseTokenName,
		BaseTokenSymbol:   baseTokenSymbol,
		BaseTokenDecimals: baseTokenDecimals,

		QuoteTokenAddress:  quoteTokenAddress,
		QuoteTokenName:     quoteTokenName,
		QuoteTokenSymbol:   quoteTokenSymbol,
		QuoteTokenDecimals: quoteTokenDecimals,

		MinOrderSize:      utils.StringToDecimal(DefaultIfNil(minOrderSize, DefaultMinOrderSize)),
		PricePrecision:    utils.ParseInt(pricePrecision, DefaultPricePrecision),
		PriceDecimals:     utils.ParseInt(priceDecimals, DefaultPriceDecimals),
		AmountDecimals:    utils.ParseInt(amountDecimals, DefaultAmountDecimals),
		MakerFeeRate:      utils.StringToDecimal(DefaultIfNil(makerFeeRate, DefaultMakerFeeRate)),
		TakerFeeRate:      utils.StringToDecimal(DefaultIfNil(takerFeeRate, DefaultTakerFeeRate)),
		GasUsedEstimation: utils.ParseInt(gasUsedEstimation, DefaultGasUsedEstimation),
	}

	err, _, _ = a.client.Post(a.MarketUrl, nil, market, nil)
	return
}

func (a *Admin) UpdateMarket(marketID, minOrderSize, pricePrecision, priceDecimals, amountDecimals, makerFeeRate, takerFeeRate, gasUsedEstimation, isPublish string) (err error) {
	fields := marketFields{
		ID:                marketID,
		MinOrderSize:      minOrderSize,
		PricePrecision:    pricePrecision,
		PriceDecimals:     priceDecimals,
		AmountDecimals:    amountDecimals,
		MakerFeeRate:      makerFeeRate,
		TakerFeeRate:      takerFeeRate,
		GasUsedEstimation: gasUsedEstimation,
		IsPublished:       isPublish,
	}

	err, _, _ = a.client.Put(a.MarketUrl, nil, fields, nil)
	return
}

func (a *Admin) PublishMarket(marketID string) (err error) {
	market := marketFields{
		ID:          marketID,
		IsPublished: "true",
	}

	err, _, _ = a.client.Put(a.MarketUrl, nil, market, nil)
	return
}

func (a *Admin) UnPublishMarket(marketID string) (err error) {
	market := marketFields{
		ID:          marketID,
		IsPublished: "false",
	}

	err, _, _ = a.client.Put(a.MarketUrl, nil, market, nil)
	return
}

func (a *Admin) UpdateMarketFee(marketID, makerFee, takerFee string) (err error) {
	market := models.Market{
		ID:           marketID,
		MakerFeeRate: utils.StringToDecimal(makerFee),
		TakerFeeRate: utils.StringToDecimal(takerFee),
	}

	err, _, _ = a.client.Put(a.MarketUrl, nil, market, nil)
	return
}

func (a *Admin) ListAccountOrders(address, limit, offset, status string) (err error) {
	var params []utils.KeyValue
	params = append(params, utils.KeyValue{Key: "address", Value: address})
	params = append(params, utils.KeyValue{Key: "limit", Value: DefaultIfNil(limit, DefaultLimit)})
	params = append(params, utils.KeyValue{Key: "offset", Value: DefaultIfNil(offset, DefaultOffset)})
	params = append(params, utils.KeyValue{Key: "status", Value: DefaultIfNil(status, DefaultStatus)})

	err, _, _ = a.client.Get(a.ListOrderUrl, params, nil, nil)
	return
}

func (a *Admin) ListAccountBalances(address, limit, offset string) (err error) {
	var params []utils.KeyValue
	params = append(params, utils.KeyValue{Key: "address", Value: address})
	params = append(params, utils.KeyValue{Key: "limit", Value: DefaultIfNil(limit, DefaultLimit)})
	params = append(params, utils.KeyValue{Key: "offset", Value: DefaultIfNil(offset, DefaultOffset)})

	err, _, _ = a.client.Get(a.ListBalanceUrl, params, nil, nil)
	return
}

func (a *Admin) ListAccountTrades(address, limit, offset, status string) (err error) {
	var params []utils.KeyValue
	params = append(params, utils.KeyValue{Key: "address", Value: address})
	params = append(params, utils.KeyValue{Key: "limit", Value: DefaultIfNil(limit, DefaultLimit)})
	params = append(params, utils.KeyValue{Key: "offset", Value: DefaultIfNil(offset, DefaultOffset)})
	params = append(params, utils.KeyValue{Key: "status", Value: DefaultIfNil(status, DefaultStatus)})

	err, _, _ = a.client.Get(a.ListTradeUrl, params, nil, nil)
	return
}

func (a *Admin) CancelOrder(ID string) (err error) {
	err, _, _ = a.client.Delete(fmt.Sprintf("%s/%s", a.CancelOrderUrl, ID), nil, nil, nil)
	return
}

func (a *Admin) RestartEngine() (err error) {
	err, _, _ = a.client.Get(a.RestartEngineUrl, nil, nil, nil)
	return
}

func DefaultIfNil(ori, dft string) string {
	if len(ori) == 0 {
		return dft
	}

	return ori
}

type marketFields struct {
	ID                string `json:"market_id"`
	MinOrderSize      string `json:"min_order_size"`
	PricePrecision    string `json:"price_precision"`
	PriceDecimals     string `json:"price_decimals"`
	AmountDecimals    string `json:"amount_decimals"`
	MakerFeeRate      string `json:"maker_fee_rate"`
	TakerFeeRate      string `json:"taker_fee_rate"`
	GasUsedEstimation string `json:"gas_used_estimation"`
	IsPublished       string `json:"is_published"`
}
