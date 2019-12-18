package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"context"
	"encoding/json"
	"github.com/HydroProtocol/nights-watch/plugin"
	"github.com/HydroProtocol/nights-watch/structs"
	"strconv"

	"github.com/HydroProtocol/hydro-scaffold-dex/backend/cli"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/connection"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/HydroProtocol/nights-watch"
	"os"
)

type DBTransactionHandler struct {
	eventQueue common.IQueue
	kvStore    common.IKVStore
}

func (handler DBTransactionHandler) TxHandlerFunc(txAndReceipt *structs.RemovableTxAndReceipt) {
	tx := txAndReceipt.Tx
	txReceipt := txAndReceipt.Receipt

	launchLog := models.LaunchLogDao.FindByHash(tx.GetHash())
	if launchLog == nil {
		utils.Debugf("Skip useless transaction %s", tx.GetHash())
		return
	}

	if launchLog.Status != common.STATUS_PENDING {
		utils.Infof("LaunchLog is not pending %s, skip", launchLog.Hash.String)
		return
	}

	txResult := txReceipt.GetResult()
	hash := tx.GetHash()
	transaction := models.TransactionDao.FindTransactionByID(launchLog.ItemID)
	utils.Infof("Transaction %s txResult is %+v", tx.GetHash(), txResult)

	var status string
	if txResult {
		status = common.STATUS_SUCCESSFUL
	} else {
		status = common.STATUS_FAILED
	}

	//approve event should not process with engine, so update and return
	if launchLog.ItemType == "hydroApprove" {
		launchLog.Status = status
		err := models.LaunchLogDao.UpdateLaunchLog(launchLog)
		if err != nil {
			panic(err)
		}
		return
	}

	event := &common.ConfirmTransactionEvent{
		Event: common.Event{
			Type:     common.EventConfirmTransaction,
			MarketID: transaction.MarketID,
		},
		Hash:      hash,
		Status:    status,
		Timestamp: txAndReceipt.TimeStamp, //todo
	}

	bts, _ := json.Marshal(event)

	err := handler.eventQueue.Push(bts)
	if err != nil {
		utils.Errorf("Push event into Queue Error: %v", err)
	}

	handler.kvStore.Set(common.HYDRO_WATCHER_BLOCK_NUMBER_CACHE_KEY, strconv.FormatUint(tx.GetBlockNumber(), 10), 0)
}

func main() {
	ctx, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	// Init Database Client
	models.Connect(os.Getenv("HSK_DATABASE_URL"))

	// Init Redis client
	client := connection.NewRedisClient(os.Getenv("HSK_REDIS_URL"))

	// init Key/Value Store
	kvStore, err := common.InitKVStore(&common.RedisKVStoreConfig{
		Ctx:    ctx,
		Client: client,
	})
	if err != nil {
		panic(err)
	}

	queue, err := common.InitQueue(&common.RedisQueueConfig{
		Name:   common.HYDRO_ENGINE_EVENTS_QUEUE_KEY,
		Client: client,
		Ctx:    ctx,
	})
	if err != nil {
		panic(err)
	}

	// only interested in tx send by launcher
	filter := func(tx sdk.Transaction) bool {
		launchLog := models.LaunchLogDao.FindByHash(tx.GetHash())

		if launchLog == nil {
			utils.Debugf("Skip useless transaction %s", tx.GetHash())
			return false
		} else {
			return true
		}
	}

	dbTxHandler := DBTransactionHandler{
		eventQueue: queue,
		kvStore:    kvStore,
	}

	p := plugin.NewTxReceiptPluginWithFilter(dbTxHandler.TxHandlerFunc, filter)

	api := os.Getenv("HSK_BLOCKCHAIN_RPC_URL")
	w := nights_watch.NewHttpBasedEthWatcher(ctx, api)
	w.RegisterTxReceiptPlugin(p)

	syncedBlockInCache, err := kvStore.Get(common.HYDRO_WATCHER_BLOCK_NUMBER_CACHE_KEY)
	if err != nil && err != common.KVStoreEmpty {
		panic(err)
	}

	var startFromBlock uint64
	if b, err := strconv.Atoi(syncedBlockInCache); err == nil {
		startFromBlock = uint64(b) + 1
	} else {
		startFromBlock = 0
	}

	go utils.StartMetrics()
	err = w.RunTillExitFromBlock(startFromBlock)
	if err != nil {
		utils.Infof("Watcher Exit with err: %s", err)
	} else {
		utils.Infof("Watcher Exit")
	}
}
