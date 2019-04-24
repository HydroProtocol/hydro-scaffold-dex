package models

import (
	"database/sql"
	"fmt"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"time"
)

type ITransactionDao interface {
	FindTransactionByHash(transactionHash string) *Transaction
	InsertTransaction(transaction *Transaction) error
	UpdateTransaction(transaction *Transaction) error
	UpdateTransactionStatus(status, hash string) error
	Count() int
	FindTransactionByID(id int64) *Transaction
}

type Transaction struct {
	ID              int64           `json:"id"              db:"id" primaryKey:"true"  autoIncrement:"true" gorm:"primary_key"`
	MarketID        string          `json:"marketID"        db:"market_id"`
	TransactionHash *sql.NullString `json:"transactionHash" db:"transaction_hash"`
	Status          string          `json:"status"          db:"status"`
	ExecutedAt      time.Time       `json:"executedAt"      db:"executed_at"`
	UpdatedAt       time.Time       `json:"updatedAt"       db:"updated_at"`
	CreatedAt       time.Time       `json:"createdAt"       db:"created_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

var TransactionDaoSqlite ITransactionDao
var TransactionDaoPG ITransactionDao

func init() {
	TransactionDaoSqlite = &transactionDaoSqlite{}
	TransactionDaoPG = &transactionDaoPG{}
}

type transactionDaoSqlite struct {
}

func (d *transactionDaoSqlite) Count() int {
	sqlString := "select count(*) from transactions"
	var count int
	err := DBSqlite.QueryRowx(sqlString).Scan(&count)

	if err != nil {
		utils.Error("GetNonce error: %v", err)
		panic(err)
	}

	return count
}

func (d *transactionDaoSqlite) FindTransactionByHash(transactionHash string) *Transaction {
	var transaction Transaction

	findBy(&transaction, &OpEq{"transaction_hash", transactionHash}, nil)

	if !transaction.TransactionHash.Valid {
		return nil
	}

	return &transaction
}

func (d *transactionDaoSqlite) InsertTransaction(transaction *Transaction) error {
	id, err := insert(transaction)

	if err != nil {
		return err
	}

	transaction.ID = id

	return nil
}

func (d *transactionDaoSqlite) FindTransactionByID(id int64) *Transaction {
	var transaction Transaction

	findBy(&transaction, &OpEq{"id", id}, nil)

	empty := Transaction{}
	if transaction == empty {
		return nil
	}

	return &transaction
}

func (*transactionDaoSqlite) UpdateTransaction(transaction *Transaction) error {
	return update(transaction, "Status", "TransactionHash", "ExecutedAt")
}

func (*transactionDaoSqlite) UpdateTransactionStatus(status, hash string) error {
	s := fmt.Sprintf(`update transactions set "status"=$1 where transaction_hash = $2`)

	_, err := DBSqlite.Exec(s, status, hash)

	return err
}

type transactionDaoPG struct {
}

func (transactionDaoPG) FindTransactionByHash(transactionHash string) *Transaction {
	var transaction Transaction
	DBPG.Where("transaction_hash = ?", transactionHash).First(&transaction)
	if !transaction.TransactionHash.Valid {
		return nil
	}

	return &transaction
}

func (transactionDaoPG) InsertTransaction(transaction *Transaction) error {
	return DBPG.Create(transaction).Error
}

func (transactionDaoPG) UpdateTransaction(transaction *Transaction) error {
	return DBPG.Save(transaction).Error
}

func (transactionDaoPG) UpdateTransactionStatus(status, hash string) error {
	return DBPG.Exec(`update transactions set "status"=$1 where transaction_hash = $2`, status, hash).Error
}

func (transactionDaoPG) Count() int {
	var count int
	DBPG.Model(&Transaction{}).Count(&count)
	return count
}

func (transactionDaoPG) FindTransactionByID(id int64) *Transaction {
	var transaction Transaction

	DBPG.Where("id = ?", id).Find(&transaction)
	if transaction.Status == "" {
		return nil
	}

	return &transaction
}
