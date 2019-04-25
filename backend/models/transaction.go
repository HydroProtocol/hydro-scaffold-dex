package models

import (
	"database/sql"
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

var TransactionDao ITransactionDao
var TransactionDaoPG ITransactionDao

func init() {
	TransactionDao = &transactionDaoPG{}
	TransactionDaoPG = TransactionDao
}

type transactionDaoPG struct {
}

func (transactionDaoPG) FindTransactionByHash(transactionHash string) *Transaction {
	var transaction Transaction
	DB.Where("transaction_hash = ?", transactionHash).First(&transaction)
	if !transaction.TransactionHash.Valid {
		return nil
	}

	return &transaction
}

func (transactionDaoPG) InsertTransaction(transaction *Transaction) error {
	return DB.Create(transaction).Error
}

func (transactionDaoPG) UpdateTransaction(transaction *Transaction) error {
	return DB.Save(transaction).Error
}

func (transactionDaoPG) UpdateTransactionStatus(status, hash string) error {
	return DB.Exec(`update transactions set "status"=$1 where transaction_hash = $2`, status, hash).Error
}

func (transactionDaoPG) Count() int {
	var count int
	DB.Model(&Transaction{}).Count(&count)
	return count
}

func (transactionDaoPG) FindTransactionByID(id int64) *Transaction {
	var transaction Transaction

	DB.Where("id = ?", id).Find(&transaction)
	if transaction.Status == "" {
		return nil
	}

	return &transaction
}
