package loud

type TradeRequest struct {
	ID               string
	Price            float64
	Amount           int
	Total            int
	IsMyTradeRequest bool
}

type ItemTradeRequest struct {
	ID               string
	TItem            Item
	Price            int
	IsMyTradeRequest bool
}

type CharacterTradeRequest struct {
	ID               string
	TCharacter       Character
	Price            int
	IsMyTradeRequest bool
}

var BuyTradeRequests = []TradeRequest{}
var SellTradeRequests = []TradeRequest{}
var SwordBuyTradeRequests = []ItemTradeRequest{}
var SwordSellTradeRequests = []ItemTradeRequest{}
var CharacterBuyTradeRequests = []CharacterTradeRequest{}
var CharacterSellTradeRequests = []CharacterTradeRequest{}
