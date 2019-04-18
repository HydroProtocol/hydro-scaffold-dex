package models

import "strings"

type ITokenDao interface {
	GetAllTokens() []*Token
	InsertToken(*Token) error
	FindTokenBySymbol(string) *Token
}

type Token struct {
	Symbol   string `json:"symbol"   db:"symbol"`
	Name     string `json:"name"     db:"name"`
	Decimals int    `json:"decimals" db:"decimals"`
	Address  string `json:"address"  db:"address"`
}

var TokenDao ITokenDao

func init() {
	TokenDao = tokenDao{}
}

type tokenDao struct {
}

func (tokenDao) InsertToken(token *Token) error {
	_, err := insert(token)
	return err
}

func (tokenDao) GetAllTokens() []*Token {
	tokens := []*Token{}
	findAllBy(&tokens, nil, nil, -1, -1)
	return tokens
}

func (tokenDao) FindTokenBySymbol(symbol string) *Token {
	var token Token
	findBy(&token, &OpEq{"symbol", symbol}, nil)

	if token.Symbol == "" {
		return nil
	}

	return &token
}

func GetBaseTokenSymbol(marketID string) string {
	splits := strings.Split(marketID, "-")

	if len(splits) != 2 {
		return ""
	} else {
		return splits[0]
	}
}

func GetBaseTokenDecimals(marketID string) int {
	tokenSymbol := GetBaseTokenSymbol(marketID)

	token := TokenDao.FindTokenBySymbol(tokenSymbol)
	if token == nil {
		panic("invalid base token, symbol:" + tokenSymbol)
	}

	return token.Decimals
}
