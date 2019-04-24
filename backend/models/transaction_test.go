package models

import (
	"database/sql"
	"github.com/HydroProtocol/hydro-sdk-backend/common"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTransactionDao_PG_UpdateTransaction(t *testing.T) {
	setEnvs()
	InitTestDBPG()

	transaction := newTransaction(common.STATUS_PENDING)
	_ = TransactionDaoPG.InsertTransaction(transaction)
	TransactionDaoPG.FindTransactionByHash(transaction.TransactionHash.String)

	dbTransaction := TransactionDaoPG.FindTransactionByID(transaction.ID)
	dbTransaction.Status = common.STATUS_SUCCESSFUL
	_ = TransactionDaoPG.UpdateTransaction(dbTransaction)

	dbTransaction2 := TransactionDaoPG.FindTransactionByID(transaction.ID)

	assert.EqualValues(t, transaction.TransactionHash, dbTransaction2.TransactionHash)
	assert.EqualValues(t, common.STATUS_SUCCESSFUL, dbTransaction2.Status)

	_ = TransactionDaoPG.UpdateTransactionStatus(common.STATUS_FAILED, dbTransaction2.TransactionHash.String)
	dbTransaction3 := TransactionDaoPG.FindTransactionByID(transaction.ID)
	assert.EqualValues(t, transaction.TransactionHash, dbTransaction3.TransactionHash)
	assert.EqualValues(t, common.STATUS_FAILED, dbTransaction3.Status)
}

func newTransaction(status string) *Transaction {
	transaction := Transaction{

		TransactionHash: &sql.NullString{uuid.NewV4().String(), true},
		MarketID:        "fix-me",
		Status:          status,
		ExecutedAt:      time.Now().UTC(),
		CreatedAt:       time.Now().UTC(),
	}

	return &transaction
}

func RandomTransaction(success bool) *Transaction {
	status := common.STATUS_SUCCESSFUL
	if !success {
		status = common.STATUS_FAILED
	}
	return newTransaction(status)
}
