package adminapi

import (
	"fmt"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/labstack/echo"
	"github.com/shopspring/decimal"
	"math/big"
	"net/http"
	"os"
	"time"
)

func RestartEngineHandler(e echo.Context) (err error) {
	restartEngineEvent := common.Event{
		Type: common.EventRestartEngine,
	}

	err = queueService.Push([]byte(utils.ToJsonString(restartEngineEvent)))
	return response(e, nil, err)
}

func GetStatusHandler(e echo.Context) (err error) {
	return response(e, map[string]interface{}{
		"web":       healthCheckService.CheckWeb(),
		"api":       healthCheckService.CheckApi(),
		"engine":    healthCheckService.CheckEngine(),
		"watcher":   healthCheckService.CheckWatcher(),
		"launcher":  healthCheckService.CheckLauncher(),
		"websocket": healthCheckService.CheckWebSocket(),
	}, err)
}

func GetBalancesHandler(e echo.Context) (err error) {
	var req struct {
		Address string `json:"address" query:"address" validate:"required"`
		Limit   int    `json:"limit"`
		Offset  int    `json:"offset"`
	}

	var resp []struct {
		Symbol        string          `json:"symbol"`
		LockedBalance decimal.Decimal `json:"lockedBalance"`
	}

	err = e.Bind(&req)
	if err == nil {
		tokens := models.TokenDao.GetAllTokens()

		for _, token := range tokens {
			lockedBalance := models.BalanceDao.GetByAccountAndSymbol(req.Address, token.Symbol, token.Decimals)
			resp = append(resp, struct {
				Symbol        string          `json:"symbol"`
				LockedBalance decimal.Decimal `json:"lockedBalance"`
			}{
				Symbol:        token.Symbol,
				LockedBalance: lockedBalance,
			},
			)
		}

		rLen := len(resp)
		if req.Offset < rLen {
			if req.Offset+req.Limit < rLen {
				resp = resp[req.Offset : req.Offset+req.Limit]
			} else {
				resp = resp[req.Offset:]
			}
		}
	}

	return response(e, map[string]interface{}{"balances": resp}, err)
}

func GetTradesHandler(e echo.Context) (err error) {
	var req struct {
		Address  string `json:"address"   query:"address"   validate:"required"`
		MarketID string `json:"market_id" query:"market_id" validate:"required"`
		Status   string `json:"status"    query:"status"`
		Offset   int    `json:"offset"    query:"offset"`
		Limit    int    `json:"limit "    query:"limit"`
	}

	var trades []*models.Trade
	var count int64
	err = e.Bind(&req)
	if err == nil {
		count, trades = models.TradeDao.FindAccountMarketTrades(req.Address, req.MarketID, req.Status, req.Offset, req.Limit)
	}

	return response(e, map[string]interface{}{"count": count, "trades": trades}, err)
}

func GetOrdersHandler(e echo.Context) (err error) {
	var req struct {
		Address  string `json:"address"   query:"address"   validate:"required"`
		MarketID string `json:"market_id" query:"market_id" validate:"required"`
		Status   string `json:"status"    query:"status"`
		Offset   int    `json:"offset"    query:"offset"`
		Limit    int    `json:"limit "    query:"limit"`
	}

	var orders []*models.Order
	var count int64

	err = e.Bind(&req)
	if err == nil {
		count, orders = models.OrderDao.FindByAccount(req.Address, req.MarketID, req.Status, req.Offset, req.Limit)
	}

	return response(e, map[string]interface{}{"count": count, "orders": orders}, err)
}

func DeleteOrderHandler(e echo.Context) (err error) {
	orderID := e.Param("order_id")

	if orderID == "" {
		err = fmt.Errorf("orderID is blank, check param")
	} else {
		order := models.OrderDao.FindByID(orderID)
		if order == nil {
			err = fmt.Errorf("cannot find order by ID %s", orderID)
		} else {
			cancelOrderEvent := common.CancelOrderEvent{
				Event: common.Event{
					Type:     common.EventCancelOrder,
					MarketID: order.MarketID,
				},
				Price: order.Price.String(),
				Side:  order.Side,
				ID:    order.ID,
			}

			err = queueService.Push([]byte(utils.ToJsonString(cancelOrderEvent)))
		}
	}

	return response(e, nil, err)
}

func ListMarketsHandler(e echo.Context) (err error) {
	markets := models.MarketDao.FindAllMarkets()
	return response(e, markets, err)
}

func EditMarketHandler(e echo.Context) (err error) {
	var fields marketFields

	err = e.Bind(&fields)
	if err != nil {
		return response(e, nil, err)
	}

	dbMarket := models.MarketDao.FindMarketByID(fields.ID)

	if dbMarket == nil {
		err = fmt.Errorf("cannot find market by ID %s", fields.ID)
		return response(e, nil, err)
	}

	var publishType string
	if dbMarket.IsPublished == false && fields.IsPublished == "true" {
		publishType = "publish"
	} else if dbMarket.IsPublished == true && fields.IsPublished == "false" {
		publishType = "unPublish"
	}

	if len(fields.MinOrderSize) > 0 {
		dbMarket.MinOrderSize = utils.StringToDecimal(fields.MinOrderSize)
	}
	if len(fields.PricePrecision) > 0 {
		dbMarket.PricePrecision = utils.ParseInt(fields.PricePrecision, 0)
	}
	if len(fields.PriceDecimals) > 0 {
		dbMarket.PriceDecimals = utils.ParseInt(fields.PriceDecimals, 0)
	}
	if len(fields.AmountDecimals) > 0 {
		dbMarket.AmountDecimals = utils.ParseInt(fields.AmountDecimals, 0)
	}
	if len(fields.MakerFeeRate) > 0 {
		dbMarket.MakerFeeRate = utils.StringToDecimal(fields.MakerFeeRate)
	}
	if len(fields.TakerFeeRate) > 0 {
		dbMarket.TakerFeeRate = utils.StringToDecimal(fields.TakerFeeRate)
	}
	if len(fields.GasUsedEstimation) > 0 {
		dbMarket.GasUsedEstimation = utils.ParseInt(fields.GasUsedEstimation, 0)
	}
	if fields.IsPublished == "true" {
		dbMarket.IsPublished = true
	} else if fields.IsPublished == "false" {
		dbMarket.IsPublished = false
	}

	err = models.MarketDao.UpdateMarket(dbMarket)
	if err == nil {
		if publishType == "publish" {
			event := common.Event{
				Type:     common.EventOpenMarket,
				MarketID: dbMarket.ID,
			}

			err = approveMarket(dbMarket)
			if err == nil {
				err = queueService.Push([]byte(utils.ToJsonString(event)))
			}
		} else if publishType == "unPublish" {
			event := common.CancelOrderEvent{
				Event: common.Event{
					Type:     common.EventCloseMarket,
					MarketID: dbMarket.ID,
				},
			}

			err = queueService.Push([]byte(utils.ToJsonString(event)))
		}
	}

	return response(e, nil, err)
}

func ApproveMarketHandler(e echo.Context) (err error) {
	marketID := e.QueryParam("marketID")
	dbMarket := models.MarketDao.FindMarketByID(marketID)

	if dbMarket == nil {
		err = fmt.Errorf("cannot find market by ID %s", marketID)
		return response(e, nil, err)
	}

	return response(e, nil, approveMarket(dbMarket))
}

func CreateMarketHandler(e echo.Context) (err error) {
	var market models.Market
	err = e.Bind(&market)
	if err != nil {
		utils.Debugf("bind param error: %v, params:%v", err, e.Request().Body)
		return response(e, nil, err)
	}

	err = models.MarketDao.InsertMarket(&market)
	return response(e, nil, err)
}

func approveMarket(market *models.Market) (err error) {
	err, quoteTokenAllowance := erc20Service.AllowanceOf(market.QuoteTokenAddress, os.Getenv("HSK_PROXY_ADDRESS"), os.Getenv("HSK_RELAYER_ADDRESS"))
	if err != nil {
		return
	}
	if quoteTokenAllowance.Cmp(big.NewInt(0)) <= 0 {
		err = approveToken(market.QuoteTokenAddress)
		if err != nil {
			return
		}
	}
	err, baseTokenAllowance := erc20Service.AllowanceOf(market.BaseTokenAddress, os.Getenv("HSK_PROXY_ADDRESS"), os.Getenv("HSK_RELAYER_ADDRESS"))
	if err != nil {
		return
	}
	if baseTokenAllowance.Cmp(big.NewInt(0)) <= 0 {
		err = approveToken(market.BaseTokenAddress)
		if err != nil {
			return
		}
	}
	return
}

func approveToken(tokenAddress string) error {
	proxyAddress := os.Getenv("HSK_PROXY_ADDRESS")
	if len(proxyAddress) != 42 {
		return fmt.Errorf("HSK_PROXY_ADDRESS empty")
	}
	proxyAddress = proxyAddress[2:]
	approveLog := models.LaunchLog{
		ItemType:  "hydroApprove",
		Status:    "created",
		From:      os.Getenv("HSK_RELAYER_ADDRESS"),
		To:        tokenAddress,
		Value:     decimal.Zero,
		GasLimit:  int64(200000),
		Data:      fmt.Sprintf("0x095ea7b3000000000000000000000000%sf000000000000000000000000000000000000000000000000000000000000000", proxyAddress),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return models.LaunchLogDao.InsertLaunchLog(&approveLog)
}

func response(e echo.Context, data interface{}, err error) error {
	ret := map[string]interface{}{}

	if err == nil {
		ret["status"] = 0
		ret["desc"] = "success"
		ret["data"] = data
	} else {
		ret["status"] = -1
		ret["desc"] = err.Error()
	}

	return e.JSONPretty(http.StatusOK, ret, "  ")
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
