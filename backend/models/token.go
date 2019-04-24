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

var TokenDao ITokenDao
var TokenDaoPG ITokenDao

func init() {
	TokenDao = &tokenDaoPG{}
	TokenDaoPG = TokenDao
}

type tokenDaoPG struct {
}

func (tokenDaoPG) GetAllTokens() []*Token {
	var tokens []*Token
	DB.Find(&tokens)
	return tokens
}

func (tokenDaoPG) InsertToken(token *Token) error {
	return DB.Create(token).Error
}

func (tokenDaoPG) FindTokenBySymbol(symbol string) *Token {
	var token Token

	DB.Where("symbol = ?", symbol).Find(&token)
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
