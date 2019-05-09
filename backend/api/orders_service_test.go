package api

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCheckBalanceAndAllowance(t *testing.T) {
	setEnvs()
	models.InitTestDBPG()
	//var marketDao models.IMarketDao
	models.MockMarketDao()
	order := BuildOrderReq{}
	address := "some address"
	checkBalanceAllowancePriceAndAmount(&order, address)
}

func TestExpiredAt(t *testing.T) {
	now := time.Now().Unix()
	timestamp := getExpiredAt(0)
	assert.EqualValues(t, true, timestamp > now)

	timestamp = getExpiredAt(5000)
	assert.True(t, timestamp > now)
}
