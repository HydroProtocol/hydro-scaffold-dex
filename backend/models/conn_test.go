package models

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	prepareTest()
	InitTestDBPG()
	var testConnect int
	_ = DBPG.Raw("select 1").Row().Scan(&testConnect)
	assert.EqualValues(t, 1, testConnect)
}

func prepareTest() {
	_ = os.Setenv("HSK_DATABASE_URL", "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	_ = os.Setenv("HSK_REDIS_URL", "redis://redis:6379/0")
	_ = os.Setenv("HSK_BLOCKCHAIN_RPC_URL", "http://127.0.0.1:8545")
	_ = os.Setenv("HSK_WETH_TOKEN_ADDRESS", "0x4a817489643a89a1428b2dd441c3fbe4dbf44789")
	_ = os.Setenv("HSK_USD_TOKEN_ADDRESS", "0xbc3524faa62d0763818636d5e400f112279d6cc0")
	_ = os.Setenv("HSK_HYDRO_TOKEN_ADDRESS", "0x4c4fa7e8ea4cfcfc93deae2c0cff142a1dd3a218")
	_ = os.Setenv("HSK_PROXY_ADDRESS", "0x04f67e8b7c39a25e100847cb167460d715215feb")
	_ = os.Setenv("HSK_HYBRID_EXCHANGE_ADDRESS", "0x179fd00c328d4ecdb5043c8686d377a24ede9d11")
	_ = os.Setenv("HSK_RELAYER_ADDRESS", "0x93388b4efe13b9b18ed480783c05462409851547")
	_ = os.Setenv("HSK_RELAYER_PK", "95b0a982c0dfc5ab70bf915dcf9f4b790544d25bc5e6cff0f38a59d0bba58651")
	_ = os.Setenv("HSK_CHAIN_ID", "50")
	_ = os.Setenv("HSK_WEB3_URL", "http://127.0.0.1:8545")
}
