package loud

type TrdReq struct {
	ID         string
	Price      float64
	Amount     int
	Total      int
	IsMyTrdReq bool
}

type ItemSellTrdReq struct {
	ID         string
	TItem      Item
	Price      int
	IsMyTrdReq bool
}

type ItemBuyTrdReq struct {
	ID         string
	TItem      ItemSpec
	Price      int
	IsMyTrdReq bool
}

type CharacterSellTrdReq struct {
	ID         string
	TCharacter Character
	Price      int
	IsMyTrdReq bool
}

type CharacterBuyTrdReq struct {
	ID         string
	TCharacter CharacterSpec
	Price      int
	IsMyTrdReq bool
}

var BuyTrdReqs = []TrdReq{}
var SellTrdReqs = []TrdReq{}
var ItemBuyTrdReqs = []ItemBuyTrdReq{}
var ItemSellTrdReqs = []ItemSellTrdReq{}
var CharacterBuyTrdReqs = []CharacterBuyTrdReq{}
var CharacterSellTrdReqs = []CharacterSellTrdReq{}
