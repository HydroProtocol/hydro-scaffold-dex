package dex_engine

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type eventHandlerSuite struct {
	suite.Suite
}

func TestRedisEventHandler(t *testing.T) {
	suite.Run(t, new(eventHandlerSuite))
}

func (s *eventHandlerSuite) TestNewOrderEvent() {
	order := newModelOrder("sell", decimal.NewFromFloat(1.0), decimal.NewFromFloat(10))

	payload, _ := json.Marshal(order)

	log.Print(string(payload))
}
