package main

import (
	"encoding/json"
	"github.com/HydroProtocol/hydro-box-dex/backend/cli"
	"github.com/HydroProtocol/hydro-box-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/launcher"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shopspring/decimal"
	"time"
)

import (
	"context"
	"os"
)

func run() int {
	_, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	//todo
	models.ConnectDatabase("sqlite3", config.Getenv("HSK_DATABASE_URL"))

	// blockchain
	hydro := ethereum.NewEthereumHydro(config.Getenv("HSK_BLOCKCHAIN_RPC_URL"))
	signService := launcher.NewDefaultSignService(config.Getenv("HSK_RELAYER_PK"), hydro.GetTransactionCount)
	gasService := func() decimal.Decimal { return utils.StringToDecimal("3000000000") } // default 10 Gwei

	launcher := launcher.NewLauncher(context.Background(), signService, hydro, gasService)

	Run(launcher, utils.StartMetrics)

	return 0
}

const pollingIntervalSeconds = 5

func Run(l *launcher.Launcher, startMetrics func()) {
	utils.Info("launcher start!")
	defer utils.Info("launcher stop!")
	go startMetrics()

	for {
		launchLogs := models.LaunchLogDao.FindAllCreated()

		if len(launchLogs) == 0 {
			select {
			case <-l.Ctx.Done():
				utils.Info("main loop Exit")
				return
			default:
				utils.Info("no logs need to be sent. sleep %ds", pollingIntervalSeconds)
				time.Sleep(pollingIntervalSeconds * time.Second)
				continue
			}
		}

		for _, launchLog := range launchLogs {
			launchLog.GasPrice = decimal.NullDecimal{
				Decimal: l.GasPrice(),
				Valid:   true,
			}

			log := launcher.LaunchLog{}
			payload, _ := json.Marshal(launchLog)

			json.Unmarshal(payload, &log)

			signedRawTransaction := l.SignService.Sign(&log)
			transactionHash, err := l.BlockChain.SendRawTransaction(signedRawTransaction)

			if err != nil {
				utils.Debug("%+v", launchLog)
				utils.Info("Send Tx failed, launchLog ID: %d, err: %+v", launchLog.ID, err)
				panic(err)
			}

			utils.Info("Send Tx, launchLog ID: %d, hash: %s", launchLog.ID, transactionHash)

			models.UpdateLaunchLogToPending(launchLog)

			if err != nil {
				utils.Info("Update Launch Log Failed, ID: %d, err: %s", launchLog.ID, err)
				panic(err)
			}

			l.SignService.AfterSign()
		}
	}
}

func main() {
	os.Exit(run())
}
