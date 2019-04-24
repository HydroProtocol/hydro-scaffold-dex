package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
)

func UpdateLaunchLogToPending(launchLog *LaunchLog) (err error) {
	launchLog.Status = common.STATUS_PENDING
	err = LaunchLogDaoSqlite.UpdateLaunchLog(launchLog)

	if err != nil {
		utils.Error("update launch error: %v", err)
		return
	}

	transaction := TransactionDaoSqlite.FindTransactionByID(launchLog.ItemID)
	transaction.TransactionHash = &launchLog.Hash

	err = TransactionDaoSqlite.UpdateTransaction(transaction)
	if err != nil {
		utils.Error("update transaction error: %v", err)
		return
	}

	trades := TradeDaoSqlite.FindTradeByTransactionID(transaction.ID)

	for _, trade := range trades {
		trade.TransactionHash = launchLog.Hash.String
		err = TradeDaoSqlite.UpdateTrade(trade)
		if err != nil {
			utils.Error("update trade error: %v", err)
			return
		}
	}
	return
}
