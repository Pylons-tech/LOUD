package loud

type Order struct {
	ID        string
	Price     float64
	Amount    int
	Total     int
	IsMyOrder bool
}

var buyOrders = []Order{}
var sellOrders = []Order{}
