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

	TokenDao.InsertToken(&token)
	tokens := TokenDao.GetAllTokens()

	assert.EqualValues(t, 1, len(tokens))
}
