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
	ID              int64           `json:"id"              db:"id" primaryKey:"true"  autoIncrement:"true"`
	MarketID        string          `json:"marketID"        db:"market_id"`
	TransactionHash *sql.NullString `json:"transactionHash" db:"transaction_hash"`
	Status          string          `json:"status"          db:"status"`
	ExecutedAt      time.Time       `json:"executedAt"      db:"executed_at"`
	UpdatedAt       time.Time       `json:"updatedAt"       db:"updated_at"`
	CreatedAt       time.Time       `json:"createdAt"       db:"created_at"`
}

var TransactionDao ITransactionDao

func init() {
	TransactionDao = &transactionDao{}
}

type transactionDao struct {
}

func (d *transactionDao) Count() int {
	sqlString := "select count(*) from transactions"
	var count int
	err := DB.QueryRowx(sqlString).Scan(&count)

	if err != nil {
		utils.Error("GetNonce error: %v", err)
		panic(err)
	}

	return count
}

func (d *transactionDao) FindTransactionByHash(transactionHash string) *Transaction {
	var transaction Transaction

	findBy(&transaction, &OpEq{"transaction_hash", transactionHash}, nil)

	if !transaction.TransactionHash.Valid {
		return nil
	}

	return &transaction
}

func (d *transactionDao) InsertTransaction(transaction *Transaction) error {
	id, err := insert(transaction)

	if err != nil {
		return err
	}

	transaction.ID = id

	return nil
}

func (d *transactionDao) FindTransactionByID(id int64) *Transaction {
	var transaction Transaction

	findBy(&transaction, &OpEq{"id", id}, nil)

	empty := Transaction{}
	if transaction == empty {
		return nil
	}

	return &transaction
}

func (*transactionDao) UpdateTransaction(transaction *Transaction) error {
	return update(transaction, "Status", "TransactionHash", "ExecutedAt")
}

func (*transactionDao) UpdateTransactionStatus(status, hash string) error {
	s := fmt.Sprintf(`update transactions set "status"=$1 where transaction_hash = $2`)

	_, err := DB.Exec(s, status, hash)

	return err
}
