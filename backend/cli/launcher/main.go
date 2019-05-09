package main

import (
	"context"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/cli"
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/launcher"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shopspring/decimal"
	"os"
	"time"
)

func run() int {
	ctx, stop := context.WithCancel(context.Background())
	go cli.WaitExitSignal(stop)

	models.Connect(os.Getenv("HSK_DATABASE_URL"))

	// blockchain
	hydro := ethereum.NewEthereumHydro(os.Getenv("HSK_BLOCKCHAIN_RPC_URL"), os.Getenv("HSK_HYBRID_EXCHANGE_ADDRESS"))
	if os.Getenv("HSK_LOG_LEVEL") == "DEBUG" {
		hydro.EnableDebug(true)
	}

	signService := launcher.NewDefaultSignService(os.Getenv("HSK_RELAYER_PK"), hydro.GetTransactionCount)

	fallbackGasPrice := decimal.New(3, 9) // 3Gwei
	priceDecider := launcher.NewGasStationGasPriceDecider(fallbackGasPrice)

	launcher := launcher.NewLauncher(ctx, signService, hydro, priceDecider)

	Run(launcher, utils.StartMetrics)

	return 0
}

const pollingIntervalSeconds = 5

func Run(l *launcher.Launcher, startMetrics func()) {
	utils.Infof("launcher start!")
	defer utils.Infof("launcher stop!")
	go startMetrics()

	for {
		launchLogs := models.LaunchLogDao.FindAllCreated()

		if len(launchLogs) == 0 {
			select {
			case <-l.Ctx.Done():
				utils.Infof("main loop Exit")
				return
			default:
				utils.Infof("no logs need to be sent. sleep %ds", pollingIntervalSeconds)

				time.Sleep(pollingIntervalSeconds * time.Second)
				continue
			}
		}

		for _, modelLaunchLog := range launchLogs {
			modelLaunchLog.GasPrice = decimal.NullDecimal{
				Decimal: l.GasPriceDecider.GasPriceInWei(),
				Valid:   true,
			}

			log := launcher.LaunchLog{
				ID:          modelLaunchLog.ID,
				ItemType:    modelLaunchLog.ItemType,
				ItemID:      modelLaunchLog.ItemID,
				Status:      modelLaunchLog.Status,
				Hash:        modelLaunchLog.Hash,
				BlockNumber: modelLaunchLog.BlockNumber,
				From:        modelLaunchLog.From,
				To:          modelLaunchLog.To,
				Value:       modelLaunchLog.Value,
				GasLimit:    modelLaunchLog.GasLimit,
				GasUsed:     modelLaunchLog.GasUsed,
				GasPrice:    modelLaunchLog.GasPrice,
				Nonce:       modelLaunchLog.Nonce,
				Data:        modelLaunchLog.Data,
				ExecutedAt:  modelLaunchLog.ExecutedAt,
				CreatedAt:   modelLaunchLog.CreatedAt,
				UpdatedAt:   modelLaunchLog.UpdatedAt,
			}
			//payload, _ := json.Marshal(launchLog)
			//json.Unmarshal(payload, &log)

			signedRawTransaction := l.SignService.Sign(&log)
			transactionHash, err := l.BlockChain.SendRawTransaction(signedRawTransaction)

			if err != nil {
				utils.Debugf("%+v", modelLaunchLog)
				utils.Infof("Send Tx failed, launchLog ID: %d, err: %+v", modelLaunchLog.ID, err)
				panic(err)
			}

			utils.Infof("Send Tx, launchLog ID: %d, hash: %s", modelLaunchLog.ID, transactionHash)

			// todo any other fields?
			modelLaunchLog.Hash = log.Hash

			models.UpdateLaunchLogToPending(modelLaunchLog)

			if err != nil {
				utils.Infof("Update Launch Log Failed, ID: %d, err: %s", modelLaunchLog.ID, err)
				panic(err)
			}

			l.SignService.AfterSign()
		}
	}
}

func main() {
	os.Exit(run())
}
