package main

import (
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/connection"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/HydroProtocol/hydro-sdk-backend/watcher"
)

func main() {
	ctx, stop := context.WithCancel(context.Background())

	go cli.WaitExitSignal(stop)

	// Init Database Client
	models.ConnectDatabase("sqlite3", config.Getenv("HSK_DATABASE_URL"))

	// Init Redis client
	client := connection.NewRedisClient(config.Getenv("HSK_REDIS_URL"))

	// Init Blockchain Client
	hydro := ethereum.NewEthereumHydro(config.Getenv("HSK_BLOCKCHAIN_RPC_URL"))

	// init Key/Value Store
	kvStore, err := common.InitKVStore(&common.RedisKVStoreConfig{
		Ctx:    ctx,
		Client: client,
	})

	if err != nil {
		panic(err)
	}

	// Init Queue
	// There is no block call of redis, so we share the client here.
	queue, err := common.InitQueue(&common.RedisQueueConfig{
		Name:   common.HYDRO_ENGINE_EVENTS_QUEUE_KEY,
		Client: client,
		Ctx:    ctx,
	})

	if err != nil {
		panic(err)
	}

	w := watcher.Watcher{
		Ctx:         ctx,
		Hydro:       hydro,
		KVClient:    kvStore,
		QueueClient: queue,
	}

	go utils.StartMetrics()
	w.Run()

	utils.Info("Watcher Exit")
}
