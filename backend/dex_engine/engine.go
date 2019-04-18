package dex_engine

import (
	"context"
	"encoding/json"
	"github.com/HydroProtocol/hydro-box-dex/backend/connection"
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/engine"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/go-redis/redis"
	"github.com/labstack/gommon/log"
	"strings"
	"sync"
)

type PgDBHandler struct {
}

func (pg PgDBHandler) Update(matchResult common.MatchResult) sync.WaitGroup {
	log.Info("testing PgDBHandler")
	return sync.WaitGroup{}
}

type RedisOrderBookSnapshotHandler struct {
	kvStore common.IKVStore
}

func (handler RedisOrderBookSnapshotHandler) Update(key string, bookSnapshot *common.SnapshotV2) sync.WaitGroup {
	bts, err := json.Marshal(bookSnapshot)
	if err != nil {
		panic(err)
	}

	_ = handler.kvStore.Set(key, string(bts), 0)

	return sync.WaitGroup{}
}

type RedisOrderBookActivitiesHandler struct {
}

func (handler RedisOrderBookActivitiesHandler) Update(webSocketMessages []common.WebSocketMessage) sync.WaitGroup {
	for _, msg := range webSocketMessages {
		if strings.HasPrefix(msg.ChannelID, "Market#") {
			pushMessage(msg)
		}
	}

	return sync.WaitGroup{}
}

type DexEngine struct {
	// all redis queues handlers
	marketHandlerMap map[string]*MarketHandler
	queue            common.IQueue

	// Wait for all queue handler exit gracefully
	Wg sync.WaitGroup

	// global ctx, if this ctx is canceled, queue handlers should exit in a short time.
	ctx context.Context

	HydroEngine *engine.Engine
	kvStore     common.IKVStore
}

func NewDexEngine(ctx context.Context, redis *redis.Client) *DexEngine {
	e := engine.NewEngine(context.Background())

	handler := PgDBHandler{}
	e.RegisterDBHandler(&handler)

	queue, _ := common.InitQueue(&common.RedisQueueConfig{
		Name:   common.HYDRO_ENGINE_EVENTS_QUEUE_KEY,
		Client: redis,
		Ctx:    ctx,
	})

	kvStore, _ := common.InitKVStore(
		&common.RedisKVStoreConfig{
			Ctx:    ctx,
			Client: redis,
		},
	)

	snapshotHandler := RedisOrderBookSnapshotHandler{kvStore: kvStore}
	e.RegisterOrderBookSnapshotHandler(snapshotHandler)

	activityHandler := RedisOrderBookActivitiesHandler{}
	e.RegisterOrderBookActivitiesHandler(activityHandler)

	engine := &DexEngine{
		queue:            queue,
		ctx:              ctx,
		marketHandlerMap: make(map[string]*MarketHandler),
		Wg:               sync.WaitGroup{},

		HydroEngine: e,
		kvStore:     kvStore,
	}

	markets := models.MarketDao.FindAllMarkets()

	for _, market := range markets {
		kvStore, _ := common.InitKVStore(
			&common.RedisKVStoreConfig{
				Ctx:    ctx,
				Client: redis,
			},
		)
		marketHandler, err := NewMarketHandler(ctx, kvStore, market, e)
		if err != nil {
			panic(err)
		}

		engine.marketHandlerMap[market.ID] = marketHandler
		utils.Info("market %s init done", marketHandler.market.ID)
	}

	return engine
}

func (e *DexEngine) start() {
	for i := range e.marketHandlerMap {
		marketHandler := e.marketHandlerMap[i]
		e.Wg.Add(1)

		go func() {
			defer e.Wg.Done()

			utils.Info("%s market handler is running", marketHandler.market.ID)
			defer utils.Info("%s market handler is stopped", marketHandler.market.ID)

			marketHandler.Run()
		}()
	}

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				for _, handler := range e.marketHandlerMap {
					close(handler.queue)
				}
				return
			default:
				data, err := e.queue.Pop()
				if err != nil {
					panic(err)
				}
				var event common.Event
				err = json.Unmarshal(data, &event)
				if err != nil {
					utils.Error("wrong event format: %+v", err)
				}

				e.marketHandlerMap[event.MarketID].queue <- data
			}
		}
	}()
}

var hydroProtocol = &ethereum.EthereumHydroProtocol{}

func Run(ctx context.Context) {
	utils.Info("dex engine start...")

	// init redis
	redisClient := connection.NewRedisClient(config.Getenv("HSK_REDIS_URL"))

	// init message queue
	messageQueue, _ := common.InitQueue(
		&common.RedisQueueConfig{
			Name:   common.HYDRO_WEBSOCKET_MESSAGES_QUEUE_KEY,
			Ctx:    ctx,
			Client: redisClient,
		},
	)
	InitWsQueue(messageQueue)

	//init database
	models.ConnectDatabase("sqlite3", config.Getenv("HSK_DATABASE_URL"))

	//start dex engine
	dexEngine := NewDexEngine(ctx, redisClient)
	dexEngine.start()

	dexEngine.Wg.Wait()
	utils.Info("dex engine stopped!")
}
