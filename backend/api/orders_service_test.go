package api

import (
	"github.com/HydroProtocol/hydro-sdk-backend/models"
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCheckBalanceAndAllowance(t *testing.T) {
	test.PreTest()
	models.InitTestDB()
	//var marketDao models.IMarketDao
	models.MockMarketDao()
	order := BuildOrderReq{}
	address := "some address"
	checkBalanceAndAllowance(&order, address)
}

func TestEetExpiredAt(t *testing.T) {
	now := time.Now().Unix()
	timestamp := getExpiredAt(0)
	assert.EqualValues(t, true, timestamp > now)

	timestamp = getExpiredAt(5000)
	assert.True(t, timestamp > now)
}
