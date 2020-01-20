package loud

type Item struct {
	ID    string `json:""`
	Name  string `json:""`
	Level int
	Price int
}

func (item *Item) GetSellPrice() int {
	switch item.Name {
	case "Wooden sword":
		if item.Level == 1 {
			return 80
		} else if item.Level == 2 {
			return 160
		}
	case "Copper sword":
		if item.Level == 1 {
			return 200
		} else if item.Level == 2 {
			return 400
		}
	}
	return -1
}
