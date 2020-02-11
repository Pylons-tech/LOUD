package loud

type Order struct {
	ID        string
	Price     float64
	Amount    int
	Total     int
	IsMyOrder bool
}

type ItemOrder struct {
	ID        string
	TItem     Item
	Price     int
	IsMyOrder bool
}

var buyOrders = []Order{}
var sellOrders = []Order{}
var swordBuyOrders = []ItemOrder{
	ItemOrder{
		"",
		Item{
			Name:  "Wooden sword",
			Level: 1,
		},
		10,
		false,
	},
	ItemOrder{
		"",
		Item{
			Name:  "Copper sword",
			Level: 1,
		},
		10,
		true,
	},
}
var swordSellOrders = []ItemOrder{
	ItemOrder{
		"",
		Item{
			Name:  "Copper sword",
			Level: 1,
		},
		10,
		false,
	},
	ItemOrder{
		"",
		Item{
			Name:  "Copper sword",
			Level: 1,
		},
		10,
		true,
	},
}
