package loud

type Order struct {
	ID     string
	Price  string
	Amount int
	Total  int
}

var buyOrders = []Order{}
var sellOrders = []Order{}