package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
)

func UpdateLaunchLogToPending(launchLog *LaunchLog) (err error) {
	launchLog.Status = common.STATUS_PENDING
	err = LaunchLogDao.UpdateLaunchLog(launchLog)

	if err != nil {
		utils.Errorf("update launch error: %v", err)
		return
	}

	//if approve event, it should not update trades or transactions
	if launchLog.ItemType == "hydroApprove" {
		return nil
	}

	transaction := TransactionDao.FindTransactionByID(launchLog.ItemID)
	transaction.TransactionHash = &launchLog.Hash

	err = TransactionDao.UpdateTransaction(transaction)
	if err != nil {
		utils.Errorf("update transaction error: %v", err)
		return
	}

	trades := TradeDao.FindTradeByTransactionID(transaction.ID)

	for _, trade := range trades {
		trade.TransactionHash = launchLog.Hash.String
		err = TradeDao.UpdateTrade(trade)
		if err != nil {
			utils.Errorf("update trade error: %v", err)
			return
		}
	}
	return
}
