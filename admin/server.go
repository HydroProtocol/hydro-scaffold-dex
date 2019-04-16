package main

import (
	"context"
	"fmt"
	"github.com/HydroProtocol/hydro-box-dex/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/connection"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"time"
)

func loadRoutes(e *echo.Echo) {
	e.Add("POST", "/markets", CreateMarketHandler)
	e.Add("PUT", "/markets", EditMarketHandler)
	e.Add("DELETE", "/orders/:orderID", DeleteOrderHandler)
	e.Add("GET", "/orders", GetOrdersHandler)
}

func GetOrdersHandler(e echo.Context) error {
	var req struct {
		Account  string `json:"account"  query:"account" validate:"required"`
		MarketID string `json:"marketID" query:"marketID" validate:"required"`
		Status   string `json:"status"   query:"status"`
		Offset   int    `json:"offset"   query:"offset"`
		limit    int    `json:"limit "   query:"limit"`
	}

	var orders []*models.Order
	var count int64
	var err error
	err = e.Bind(req)
	if err == nil {
		count, orders = models.OrderDao.FindByAccount(req.Account, req.MarketID, req.Status, req.Offset, req.limit)
	}

	return response(e, map[string]interface{}{"count": count, "orders": orders}, err)
}

func DeleteOrderHandler(e echo.Context) error {
	orderID := e.Param("orderID")
	var err error
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

			err = QueueService.Push([]byte(utils.ToJsonString(cancelOrderEvent)))
		}
	}

	return response(e, nil, err)
}

func EditMarketHandler(e echo.Context) error {
	market := &models.Market{}
	e.Bind(market)
	dbMarket := models.MarketDao.FindMarketByID(market.ID)
	var err error
	if dbMarket != nil {
		err = fmt.Errorf("cannot find market by ID %s", market.ID)
	} else {
		err = models.MarketDao.UpdateMarket(market)
	}

	return response(e, nil, err)
}

func CreateMarketHandler(e echo.Context) error {
	market := &models.Market{}
	e.Bind(market)

	err := models.MarketDao.InsertMarket(market)
	return response(e, nil, err)
}

func response(e echo.Context, data interface{}, err error) error {
	ret := map[string]interface{}{}

	if err == nil {
		ret["status"] = "success"
		ret["data"] = data
	} else {
		ret["status"] = "fail"
		ret["error_message"] = err.Error()
	}

	return e.JSON(http.StatusOK, ret)
}

func newEchoServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	loadRoutes(e)

	return e
}

func StartServer(ctx context.Context) {
	//init database
	models.ConnectDatabase("sqlite3", config.Getenv("HSK_DATABASE_URL"))

	//init event queue
	QueueService, _ = common.InitQueue(
		&common.RedisQueueConfig{
			Name:   common.HYDRO_ENGINE_EVENTS_QUEUE_KEY,
			Ctx:    ctx,
			Client: connection.NewRedisClient(config.Getenv("HSK_REDIS_URL")),
		},
	)

	e := newEchoServer()
	s := &http.Server{
		Addr:         ":3003",
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	go func() {
		if err := e.StartServer(s); err != nil {
			e.Logger.Info("shutting down the server: %v", err)
			panic(err)
		}
	}()

	<-ctx.Done()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
