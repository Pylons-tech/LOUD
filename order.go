package loud

type TradeRequest struct {
	ID        string
	Price     float64
	Amount    int
	Total     int
	IsMyTradeRequest bool
}

type ItemTradeRequest struct {
	ID        string
	TItem     Item
	Price     int
	IsMyTradeRequest bool
}

var buyTradeRequests = []TradeRequest{}
var sellTradeRequests = []TradeRequest{}
var swordBuyTradeRequests = []ItemTradeRequest{}
var swordSellTradeRequests = []ItemTradeRequest{}
