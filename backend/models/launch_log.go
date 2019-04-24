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
	ID       int64          `db:"id" auto:"true" primaryKey:"true" autoIncrement:"true" gorm:"primary_key"`
	ItemType string         `db:"item_type"`
	ItemID   int64          `db:"item_id"`
	Status   string         `db:"status"`
	Hash     sql.NullString `db:"transaction_hash" gorm:"column:transaction_hash"`

	BlockNumber sql.NullInt64 `db:"block_number"`

	From     string              `db:"t_from" gorm:"column:t_from"`
	To       string              `db:"t_to" gorm:"column:t_to"`
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

var LaunchLogDaoSqlite ILaunchLogDao
var LaunchLogDaoPG ILaunchLogDao

func init() {
	LaunchLogDaoSqlite = &launchLogDaoSqlite{}
	LaunchLogDaoPG = &launchLogDaoPG{}
}

type launchLogDaoSqlite struct {
}

func (launchLogDaoSqlite) UpdateLaunchLogsStatusByItemID(status string, itemID int64) error {
	_, err := DBSqlite.Exec(`update launch_logs set "status" = $1 where item_id = $2`, status, itemID)
	return err
}

func (launchLogDaoSqlite) FindLaunchLogByID(id int) *LaunchLog {
	var launchLog LaunchLog
	findBy(&launchLog, &OpEq{"id", id}, nil)

	return &launchLog
}

func (launchLogDaoSqlite) FindByHash(hash string) *LaunchLog {
	var launchLog LaunchLog

	findBy(&launchLog, &OpEq{"transaction_hash", hash}, nil)

	if !launchLog.Hash.Valid {
		return nil
	}

	return &launchLog
}

func (launchLogDaoSqlite) FindPendingLogWithMaxNonce() int64 {
	var nonce sql.NullInt64
	err := DBSqlite.QueryRow(`select max(nonce) from launch_logs`).Scan(&nonce)
	if err != nil {
		panic(err)
	}

	if nonce.Valid {
		return nonce.Int64
	} else {
		return -1
	}
}

func (launchLogDaoSqlite) FindAllCreated() []*LaunchLog {
	var launchLogs []*LaunchLog

	findAllBy(&launchLogs,
		&OpEq{"status", "created"},
		map[string]OrderByDirection{"created_at": OrderByAsc}, -1, -1)

	return launchLogs
}

func (launchLogDaoSqlite) UpdateLaunchLog(launchLog *LaunchLog) error {
	return update(launchLog, "ItemID", "Status", "Hash", "BlockNumber", "From", "To", "Value", "GasLimit", "GasUsed", "GasPrice", "Nonce", "Data", "ExecutedAt", "CreatedAt", "UpdatedAt")
}

func (launchLogDaoSqlite) InsertLaunchLog(launchLog *LaunchLog) error {
	id, err := insert(launchLog)

	if err != nil {
		return err
	}

	launchLog.ID = id

	return nil
}

//pg

type launchLogDaoPG struct {
}

func (launchLogDaoPG) FindLaunchLogByID(id int) *LaunchLog {
	var launchLog LaunchLog

	DBPG.First(&launchLog, id)
	return &launchLog
}

func (launchLogDaoPG) FindByHash(hash string) *LaunchLog {
	var launchLog LaunchLog

	DBPG.Where("transaction_hash = ?", hash).Find(&launchLog)
	if !launchLog.Hash.Valid {
		return nil
	}

	return &launchLog
}

func (launchLogDaoPG) FindPendingLogWithMaxNonce() int64 {
	var nonce sql.NullInt64

	err := DBPG.Raw(`select max(nonce) from launch_logs`).Row().Scan(&nonce)
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
	DBPG.Where("status = 'created'").Order("created_at asc").Find(&launchLogs)
	return launchLogs
}

func (launchLogDaoPG) UpdateLaunchLog(launchLog *LaunchLog) error {
	return DBPG.Save(launchLog).Error
}

func (launchLogDaoPG) InsertLaunchLog(launchLog *LaunchLog) error {
	return DBPG.Create(launchLog).Error
}

func (launchLogDaoPG) UpdateLaunchLogsStatusByItemID(status string, itemID int64) error {
	return DBPG.Exec(`update launch_logs set "status" = ? where item_id = ?`, status, itemID).Error
}
