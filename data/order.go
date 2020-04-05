package loud

type TradeRequest struct {
	ID               string
	Price            float64
	Amount           int
	Total            int
	IsMyTradeRequest bool
}

type ItemSellTradeRequest struct {
	ID               string
	TItem            Item
	Price            int
	IsMyTradeRequest bool
}

type ItemBuyTradeRequest struct {
	ID               string
	TItem            ItemSpec
	Price            int
	IsMyTradeRequest bool
}

type CharacterSellTradeRequest struct {
	ID               string
	TCharacter       Character
	Price            int
	IsMyTradeRequest bool
}

type CharacterBuyTradeRequest struct {
	ID               string
	TCharacter       CharacterSpec
	Price            int
	IsMyTradeRequest bool
}

var BuyTradeRequests = []TradeRequest{}
var SellTradeRequests = []TradeRequest{}
var SwordBuyTradeRequests = []ItemBuyTradeRequest{}
var SwordSellTradeRequests = []ItemSellTradeRequest{}
var CharacterBuyTradeRequests = []CharacterBuyTradeRequest{}
var CharacterSellTradeRequests = []CharacterSellTradeRequest{}
