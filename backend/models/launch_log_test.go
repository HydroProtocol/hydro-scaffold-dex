package models

import (
	"database/sql"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"github.com/davecgh/go-spew/spew"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestLaunchLogDao_PG_FindAllCreated(t *testing.T) {
	setEnvs()
	InitTestDBPG()

	launchLogs := LaunchLogDaoPG.FindAllCreated()
	assert.EqualValues(t, len(launchLogs), 0)

	launchLog1 := newLaunchLog()
	launchLog2 := newLaunchLog()
	launchLog3 := newLaunchLog()

	spew.Dump(launchLog1)
	spew.Dump(LaunchLogDaoPG.InsertLaunchLog(launchLog1))
	//_ = LaunchLogDaoPG.InsertLaunchLog(launchLog1)
	_ = LaunchLogDaoPG.InsertLaunchLog(launchLog2)
	_ = LaunchLogDaoPG.InsertLaunchLog(launchLog3)

	launchLogs = LaunchLogDaoPG.FindAllCreated()
	assert.EqualValues(t, 3, len(launchLogs))

	launchLog := LaunchLogDaoPG.FindLaunchLogByID(1)
	assert.EqualValues(t, 1, launchLog.ID)
	assert.EqualValues(t, "created", launchLog.Status)
	launchLog.Status = common.STATUS_PENDING
	_ = LaunchLogDaoPG.UpdateLaunchLog(launchLog)

	launchLog = LaunchLogDaoPG.FindLaunchLogByID(1)
	assert.EqualValues(t, 1, launchLog.ID)
	assert.EqualValues(t, common.STATUS_PENDING, launchLog.Status)
}

func newLaunchLog() *LaunchLog {
	launchLog := LaunchLog{
		ItemType:    "hydro_trade",
		ItemID:      int64(rand.Int31()),
		Status:      "created",
		Hash:        sql.NullString{},
		BlockNumber: sql.NullInt64{},

		From:     TestUser1,
		To:       TestUser2,
		Value:    decimal.Zero,
		GasLimit: 1000000,
		GasPrice: decimal.NullDecimal{Decimal: utils.StringToDecimal("1000000000"), Valid: true},
		Nonce:    sql.NullInt64{},
		Data:     "some data",

		ExecutedAt: time.Now().UTC(),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	return &launchLog
}
