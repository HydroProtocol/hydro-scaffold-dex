package models

import (
	"database/sql"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/config"
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestLaunchLogDao_FindAllCreated(t *testing.T) {
	test.PreTest()
	InitTestDB()

	launchLogs := LaunchLogDao.FindAllCreated()
	assert.EqualValues(t, len(launchLogs), 0)

	launchLog1 := newLaunchLog()
	launchLog2 := newLaunchLog()
	launchLog3 := newLaunchLog()

	_ = LaunchLogDao.InsertLaunchLog(launchLog1)
	_ = LaunchLogDao.InsertLaunchLog(launchLog2)
	_ = LaunchLogDao.InsertLaunchLog(launchLog3)

	launchLogs = LaunchLogDao.FindAllCreated()
	assert.EqualValues(t, 3, len(launchLogs))

	launchLog := LaunchLogDao.FindLaunchLogByID(1)
	assert.EqualValues(t, 1, launchLog.ID)
	assert.EqualValues(t, "created", launchLog.Status)
	launchLog.Status = common.STATUS_PENDING
	_ = LaunchLogDao.UpdateLaunchLog(launchLog)

	launchLog = LaunchLogDao.FindLaunchLogByID(1)
	assert.EqualValues(t, 1, launchLog.ID)
	assert.EqualValues(t, common.STATUS_PENDING, launchLog.Status)
}

func newLaunchLog() *LaunchLog {
	launchLog := LaunchLog{
		ItemType:    "hydro_trade",
		ItemID:      rand.Int63(),
		Status:      "created",
		Hash:        sql.NullString{},
		BlockNumber: sql.NullInt64{},

		From:     config.User1,
		To:       config.User2,
		Value:    decimal.Zero,
		GasLimit: 1000000,
		GasPrice: decimal.NullDecimal{utils.StringToDecimal("1000000000"), true},
		Nonce:    sql.NullInt64{},
		Data:     "some data",

		ExecutedAt: time.Now(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return &launchLog
}
