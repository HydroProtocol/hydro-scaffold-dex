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
type LaunchLog struct {
	ID          int64          `db:"id" auto:"true" primaryKey:"true" autoIncrement:"true" gorm:"primary_key"`
	ItemType    string         `db:"item_type"`
	ItemID      int64          `db:"item_id"`
	Status      string         `db:"status"`
	Hash        sql.NullString `db:"transaction_hash" gorm:"column:transaction_hash"`
	BlockNumber sql.NullInt64  `db:"block_number"`

	From     string              `db:"t_from" gorm:"column:t_from"`
	To       string              `db:"t_to"   gorm:"column:t_to"`
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

func (LaunchLog) TableName() string {
	return "launch_logs"
}

var LaunchLogDao ILaunchLogDao
var LaunchLogDaoPG ILaunchLogDao

func init() {
	LaunchLogDao = &launchLogDaoPG{}
	LaunchLogDaoPG = LaunchLogDao
}

type launchLogDaoPG struct {
}

func (launchLogDaoPG) FindLaunchLogByID(id int) *LaunchLog {
	var launchLog LaunchLog

	DB.First(&launchLog, id)
	return &launchLog
}

func (launchLogDaoPG) FindByHash(hash string) *LaunchLog {
	var launchLog LaunchLog

	DB.Where("transaction_hash = ?", hash).Find(&launchLog)
	if !launchLog.Hash.Valid {
		return nil
	}

	return &launchLog
}

func (launchLogDaoPG) FindPendingLogWithMaxNonce() int64 {
	var nonce sql.NullInt64

	err := DB.Raw(`select max(nonce) from launch_logs`).Row().Scan(&nonce)
	if err != nil {
		panic(err)
	}
	if nonce.Valid {
		return nonce.Int64
	} else {
		return -1
	}
}

func (launchLogDaoPG) FindAllCreated() []*LaunchLog {
	var launchLogs []*LaunchLog
	DB.Where("status = 'created'").Order("created_at asc").Find(&launchLogs)
	return launchLogs
}

func (launchLogDaoPG) UpdateLaunchLog(launchLog *LaunchLog) error {
	return DB.Save(launchLog).Error
}

func (launchLogDaoPG) InsertLaunchLog(launchLog *LaunchLog) error {
	return DB.Create(launchLog).Error
}

func (launchLogDaoPG) UpdateLaunchLogsStatusByItemID(status string, itemID int64) error {
	return DB.Exec(`update launch_logs set "status" = ? where item_id = ?`, status, itemID).Error
}
