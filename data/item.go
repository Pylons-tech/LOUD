package loud

type Item struct {
	ID    string `json:""`
	Name  string `json:""`
	Level int
	Price int
}

const (
	WOODEN_SWORD string = "Wooden sword"
	COPPER_SWORD        = "Copper sword"
)

var ShopItems = []Item{
	Item{
		ID:    "001",
		Name:  WOODEN_SWORD,
		Level: 1,
		Price: 100,
	},
	Item{
		ID:    "002",
		Name:  COPPER_SWORD,
		Level: 1,
		Price: 250,
	},
}

var WorldItems = []Item{
	Item{
		Name:  WOODEN_SWORD,
		Level: 1,
	},
	Item{
		Name:  WOODEN_SWORD,
		Level: 2,
	},
	Item{
		Name:  COPPER_SWORD,
		Level: 1,
	},
	Item{
		Name:  COPPER_SWORD,
		Level: 2,
	},
}

func (item *Item) GetSellPrice() int {
	switch item.Name {
	case WOODEN_SWORD:
		if item.Level == 1 {
			return 80
		} else if item.Level == 2 {
			return 160
		}
	case COPPER_SWORD:
		if item.Level == 1 {
			return 200
		} else if item.Level == 2 {
			return 400
		}
	}
	return -1
}

func (item *Item) GetUpgradePrice() int {
	switch item.Name {
	case WOODEN_SWORD:
		return 250
	case COPPER_SWORD:
		return 100
	}
	return -1
}
