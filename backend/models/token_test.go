package models

import (
	"github.com/HydroProtocol/hydro-sdk-backend/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenDao_GetAllTokens(t *testing.T) {
	test.PreTest()
	InitTestDB()

	token := Token{
		Address:  "some address",
		Name:     "HOT",
		Symbol:   "HOT",
		Decimals: 18,
	}

	TokenDaoSqlite.InsertToken(&token)
	tokens := TokenDaoSqlite.GetAllTokens()

	assert.EqualValues(t, 1, len(tokens))
}

//pg
func TestTokenDao_PG_GetAllTokens(t *testing.T) {
	prepareTest()
	InitTestDBPG()

	token := Token{
		Address:  "some address",
		Name:     "HOT",
		Symbol:   "HOT",
		Decimals: 18,
	}

	TokenDaoPG.InsertToken(&token)
	tokens := TokenDaoPG.GetAllTokens()

	assert.EqualValues(t, 1, len(tokens))
}
