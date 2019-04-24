package models

import "strings"

type ITokenDao interface {
	GetAllTokens() []*Token
	InsertToken(*Token) error
	FindTokenBySymbol(string) *Token
}

type Token struct {
	Symbol   string `json:"symbol"   db:"symbol" gorm:"primary_key"`
	Name     string `json:"name"     db:"name"`
	Decimals int    `json:"decimals" db:"decimals"`
	Address  string `json:"address"  db:"address"`
}

func (Token) TableName() string {
	return "tokens"
}

var TokenDaoSqlite ITokenDao
var TokenDaoPG ITokenDao

func init() {
	TokenDaoSqlite = tokenDaoSqlite{}
	TokenDaoPG = tokenDaoPG{}
}

type tokenDaoSqlite struct {
}

func (tokenDaoSqlite) InsertToken(token *Token) error {
	_, err := insert(token)
	return err
}

func (tokenDaoSqlite) GetAllTokens() []*Token {
	tokens := []*Token{}
	findAllBy(&tokens, nil, nil, -1, -1)
	return tokens
}

func (tokenDaoSqlite) FindTokenBySymbol(symbol string) *Token {
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

	token := TokenDaoSqlite.FindTokenBySymbol(tokenSymbol)
	if token == nil {
		panic("invalid base token, symbol:" + tokenSymbol)
	}

	return token.Decimals
}

//pg
type tokenDaoPG struct {
}

func (tokenDaoPG) GetAllTokens() []*Token {
	var tokens []*Token
	DBPG.Find(&tokens)
	return tokens
}

func (tokenDaoPG) InsertToken(token *Token) error {
	return DBPG.Create(token).Error
}

func (tokenDaoPG) FindTokenBySymbol(symbol string) *Token {
	var token Token

	DBPG.Where("symbol = ?", symbol).Find(&token)
	if token.Symbol == "" {
		return nil
	}

	return &token
}
