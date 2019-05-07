package api

import (
	"github.com/HydroProtocol/hydro-scaffold-dex/backend/models"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

const MaxBarsCount = 200

func GetAllTrades(p Param) (interface{}, error) {
	req := p.(*QueryTradeReq)
	count, trades := models.TradeDao.FindAllTrades(req.MarketID)

	resp := QueryTradeResp{
		Count:  count,
		Trades: trades,
	}
	return &resp, nil
}

func GetAccountTrades(p Param) (interface{}, error) {
	req := p.(*QueryTradeReq)
	if req.PerPage <= 0 {
		req.PerPage = 20
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	offset := req.PerPage * (req.Page - 1)
	limit := req.PerPage

	count, trades := models.TradeDao.FindAccountMarketTrades(req.Address, req.MarketID, req.Status, limit, offset)

	return &QueryTradeResp{
		Count:  count,
		Trades: trades,
	}, nil
}

func GetTradingView(p Param) (interface{}, error) {
	params := p.(*CandlesReq)
	pair := params.MarketID
	from := params.From
	to := params.To
	granularity := params.Granularity

	if (to - granularity*MaxBarsCount) > from {
		from = to - granularity*MaxBarsCount
	}

	trades := models.TradeDao.FindTradesByMarket(pair, time.Unix(from, 0), time.Unix(to, 0))

	if len(trades) == 0 {
		return map[string]interface{}{
			"candles": []*Bar{},
		}, nil
	}

	return map[string]interface{}{
		"candles": BuildTradingViewByTrades(trades, granularity),
	}, nil

}

func BuildTradingViewByTrades(trades []*models.Trade, granularity int64) []*Bar {
	var bars []*Bar
	var currentIndex int64
	var currentBar *Bar

	sort.Slice(trades, func(i, j int) bool {
		return trades[i].ExecutedAt.Unix() < trades[j].ExecutedAt.Unix()
	})

	for _, trade := range trades {
		tIndex := trade.ExecutedAt.Unix() / granularity
		if currentBar == nil || currentBar.Volume.IsZero() {
			currentIndex = tIndex
			currentBar = newBar(trade, currentIndex, granularity)
			continue
		}

		if tIndex < currentIndex+1 {
			currentBar.High = decimal.Max(currentBar.High, trade.Price)
			currentBar.Low = decimal.Min(currentBar.Low, trade.Price)
			currentBar.Volume = currentBar.Volume.Add(trade.Amount)
			currentBar.Close = trade.Price
		} else {
			currentIndex = tIndex
			if currentBar.Volume.IsZero() {
				continue
			}
			bars = pushBar(bars, currentBar)
			currentBar = newBar(trade, currentIndex, granularity)
		}
	}

	bars = pushBar(bars, currentBar)

	return bars
}

func pushBar(bars []*Bar, bar *Bar) []*Bar {
	newBar := &Bar{
		Time:   bar.Time,
		Open:   bar.Open,
		Close:  bar.Close,
		Low:    bar.Low,
		High:   bar.High,
		Volume: bar.Volume,
	}

	bars = append(bars, newBar)
	return bars
}

func newBar(trade *models.Trade, currentIndex int64, granularity int64) *Bar {
	bar := &Bar{
		Time:   currentIndex * granularity,
		Volume: trade.Amount,
		Open:   trade.Price,
		Close:  trade.Price,
		High:   trade.Price,
		Low:    trade.Price,
	}

	return bar
}

type Bar struct {
	Time   int64           `json:"time"`
	Open   decimal.Decimal `json:"open"`
	Close  decimal.Decimal `json:"close"`
	Low    decimal.Decimal `json:"low"`
	High   decimal.Decimal `json:"high"`
	Volume decimal.Decimal `json:"volume"`
}
