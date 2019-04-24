package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenDao_PG_GetAllTokens(t *testing.T) {
	setEnvs()
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
