package models

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"time"
)

type ILaunchLogDao interface {
	FindLaunchLogByID(int) *LaunchLog
	FindByHash(hash string) *LaunchLog
	FindPendingLogWithMaxNonce() int64
	FindAllCreated() []*LaunchLog
	UpdateLaunchLog(*LaunchLog) error
	InsertLaunchLog(*LaunchLog) error
	UpdateLaunchLogsStatusByItemID(string, int64) error
}

var LaunchLogDao ILaunchLogDao

func init() {
	LaunchLogDao = &launchLogDao{}
}

type launchLogDao struct {
}

func (launchLogDao) UpdateLaunchLogsStatusByItemID(status string, itemID int64) error {
	_, err := DB.Exec(`update launch_logs set "status" = $1 where item_id = $2`, status, itemID)
	return err
}

func (launchLogDao) FindLaunchLogByID(id int) *LaunchLog {
	var launchLog LaunchLog
	findBy(&launchLog, &OpEq{"id", id}, nil)

	return &launchLog
}

func (launchLogDao) FindByHash(hash string) *LaunchLog {
	var launchLog LaunchLog

	findBy(&launchLog, &OpEq{"transaction_hash", hash}, nil)

	if !launchLog.Hash.Valid {
		return nil
	}

	return &launchLog
}

func (launchLogDao) FindPendingLogWithMaxNonce() int64 {
	var nonce sql.NullInt64
	err := DB.QueryRow(`select max(nonce) from launch_logs`).Scan(&nonce)
	if err != nil {
		panic(err)
	}

	if nonce.Valid {
		return nonce.Int64
	} else {
		return -1
	}
}

func (launchLogDao) FindAllCreated() []*LaunchLog {
	launchLogs := []*LaunchLog{}
	findAllBy(&launchLogs,
		&OpEq{"status", "created"},
		map[string]OrderByDirection{"created_at": OrderByAsc}, -1, -1)

	return launchLogs
}

func (launchLogDao) UpdateLaunchLog(launchLog *LaunchLog) error {
	return update(launchLog, "ItemID", "Status", "Hash", "BlockNumber", "From", "To", "Value", "GasLimit", "GasUsed", "GasPrice", "Nonce", "Data", "ExecutedAt", "CreatedAt", "UpdatedAt")
}

func (launchLogDao) InsertLaunchLog(launchLog *LaunchLog) error {
	id, err := insert(launchLog)

	if err != nil {
		return err
	}

	launchLog.ID = id

	return nil
}

func init() {
	LaunchLogDao = launchLogDao{}
}

type LaunchLog struct {
	ID       int64          `db:"id" auto:"true" primaryKey:"true" autoIncrement:"true"`
	ItemType string         `db:"item_type"`
	ItemID   int64          `db:"item_id"`
	Status   string         `db:"status"`
	Hash     sql.NullString `db:"transaction_hash"`

	BlockNumber sql.NullInt64 `db:"block_number"`

	From     string              `db:"t_from"`
	To       string              `db:"t_to"`
	Value    decimal.Decimal     `db:"value"`
	GasLimit int64               `db:"gas_limit"`
	GasUsed  sql.NullInt64       `db:"gas_used"`
	GasPrice decimal.NullDecimal `db:"gas_price"`
	Nonce    sql.NullInt64       `db:"nonce"`
	Data     string              `db:"data"`

	ExecutedAt time.Time `db:"executed_at"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
