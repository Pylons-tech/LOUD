package loud

// TrdReq is a struct to manage gold buy/sell trade requests
type TrdReq struct {
	ID         string
	Price      float64
	Amount     int
	Total      int
	IsMyTrdReq bool
}

// ItemSellTrdReq is a struct to manage item sell trade requests
type ItemSellTrdReq struct {
	ID         string
	TItem      Item
	Price      int
	IsMyTrdReq bool
}

// ItemBuyTrdReq is a struct to manage item buy trade requests
type ItemBuyTrdReq struct {
	ID         string
	TItem      ItemSpec
	Price      int
	IsMyTrdReq bool
}

// CharacterSellTrdReq is a struct to manage character sell trade requests
type CharacterSellTrdReq struct {
	ID         string
	TCharacter Character
	Price      int
	IsMyTrdReq bool
}

// CharacterBuyTrdReq is a struct to manage character buy trade requests
type CharacterBuyTrdReq struct {
	ID         string
	TCharacter CharacterSpec
	Price      int
	IsMyTrdReq bool
}

// BuyTrdReqs is a global variable to store gold buy trade requests
var BuyTrdReqs = []TrdReq{}

// SellTrdReqs is a global variable to store gold sell trade requests
var SellTrdReqs = []TrdReq{}

// ItemBuyTrdReqs is a global variable to store item buy trade requests
var ItemBuyTrdReqs = []ItemBuyTrdReq{}

// ItemSellTrdReqs is a global variable to store item sell trade requests
var ItemSellTrdReqs = []ItemSellTrdReq{}

// CharacterBuyTrdReqs is a global variable to store character buy trade requests
var CharacterBuyTrdReqs = []CharacterBuyTrdReq{}

// CharacterSellTrdReqs is a global variable to store character sell trade requests
var CharacterSellTrdReqs = []CharacterSellTrdReq{}
