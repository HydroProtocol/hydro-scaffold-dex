package main

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

import (
	"context"
	"encoding/json"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/connection"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/HydroProtocol/hydro-sdk-backend/watcher"
)

type DBTransactionHandler struct {
	w watcher.Watcher
}

func (handler DBTransactionHandler) Update(tx sdk.Transaction, timestamp uint64) {
	launchLog := models.LaunchLogDao.FindByHash(tx.GetHash())

	if launchLog == nil {
		utils.Debug("Skip useless transaction %s", tx.GetHash())
		return
	}

	if launchLog.Status != common.STATUS_PENDING {
		utils.Info("LaunchLog is not pending %s, skip", launchLog.Hash.String)
		return
	}

	if launchLog != nil {
		txReceipt, _ := handler.w.Hydro.GetTransactionReceipt(tx.GetHash())
		result := txReceipt.GetResult()
		hash := tx.GetHash()
		transaction := models.TransactionDao.FindTransactionByID(launchLog.ItemID)
		utils.Info("Transaction %s result is %+v", tx.GetHash(), result)
		//w.handleTransaction(launchLog.ItemID, result)

		var status string

		if result {
			status = common.STATUS_SUCCESSFUL
		} else {
			status = common.STATUS_FAILED
		}

		event := &common.ConfirmTransactionEvent{
			Event: common.Event{
				Type:     common.EventConfirmTransaction,
				MarketID: transaction.MarketID,
			},
			Hash:      hash,
			Status:    status,
			Timestamp: timestamp,
		}

		bts, _ := json.Marshal(event)

		err := handler.w.QueueClient.Push(bts)

		if err != nil {
			utils.Error("Push event into Queue Error %v", err)
		}
	}
}

func main() {
	ctx, stop := context.WithCancel(context.Background())

	go cli.WaitExitSignal(stop)

	// Init Database Client
	models.Connect(config.Getenv("HSK_DATABASE_URL"))

	// Init Redis client
	client := connection.NewRedisClient(os.Getenv("HSK_REDIS_URL"))

	// Init Blockchain Client
	hydro := ethereum.NewEthereumHydro(os.Getenv("HSK_BLOCKCHAIN_RPC_URL"), os.Getenv("HSK_HYBRID_EXCHANGE_ADDRESS"))
	if os.Getenv("HSK_LOG_LEVEL") == "DEBUG" {
		hydro.EnableDebug(true)
	}

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

	w.RegisterHandler(DBTransactionHandler{w})

	go utils.StartMetrics()

	w.Run()

	utils.Info("Watcher Exit")
}
