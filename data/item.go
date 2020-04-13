package loud

import "strings"

type Item struct {
	ID         string `json:""`
	Name       string `json:""`
	Level      int
	Attack     int
	Price      int
	PreItem    string `json:""`
	LastUpdate int64
}

type ItemSpec struct {
	Name   string `json:""`
	Level  [2]int
	Attack [2]int
	Price  int
}

type Character struct {
	ID         string `json:""`
	Name       string `json:""`
	Level      int
	Price      int
	XP         float64
	HP         int
	MaxHP      int
	LastUpdate int64
}
type CharacterSpec struct {
	Name  string `json:""`
	Level [2]int
	Price int
	XP    [2]float64
	HP    [2]int
	MaxHP [2]int
}

const (
	WOODEN_SWORD string = "Wooden sword"
	COPPER_SWORD        = "Copper sword"
	SILVER_SWORD        = "Silver sword"
	BRONZE_SWORD        = "Bronze sword"
	IRON_SWORD          = "Iron sword"
	GOBLIN_EAR          = "Goblin ear"
	WOLF_TAIL           = "Wolf tail"
	TROLL_TOES          = "Troll toes"
)

func (item Item) IsSword() bool {
	return strings.Contains(item.Name, "sword")
}

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
	Item{
		ID:      "003",
		Name:    SILVER_SWORD,
		Level:   1,
		Price:   250,
		PreItem: GOBLIN_EAR,
	},
	Item{
		ID:      "004",
		Name:    BRONZE_SWORD,
		Level:   1,
		Price:   250,
		PreItem: WOLF_TAIL,
	},
	Item{
		ID:      "005",
		Name:    IRON_SWORD,
		Level:   1,
		Price:   250,
		PreItem: TROLL_TOES,
	},
}

var WorldItemSpecs = []ItemSpec{
	ItemSpec{
		Name:   WOODEN_SWORD,
		Level:  [2]int{1, 1},
		Attack: [2]int{3, 3},
	},
	ItemSpec{
		Name:   WOODEN_SWORD,
		Level:  [2]int{2, 2},
		Attack: [2]int{6, 6},
	},
	ItemSpec{
		Name:   COPPER_SWORD,
		Level:  [2]int{1, 1},
		Attack: [2]int{10, 10},
	},
	ItemSpec{
		Name:   COPPER_SWORD,
		Level:  [2]int{2, 2},
		Attack: [2]int{20, 20},
	},
	ItemSpec{
		Name:  TROLL_TOES,
		Level: [2]int{1, 1},
	},
	ItemSpec{
		Name:  WOLF_TAIL,
		Level: [2]int{1, 1},
	},
	ItemSpec{
		Name:  GOBLIN_EAR,
		Level: [2]int{1, 1},
	},
	ItemSpec{
		Name:   SILVER_SWORD,
		Level:  [2]int{1, 1},
		Attack: [2]int{30, 30},
	},
	ItemSpec{
		Name:   BRONZE_SWORD,
		Level:  [2]int{1, 1},
		Attack: [2]int{50, 50},
	},
	ItemSpec{
		Name:   IRON_SWORD,
		Level:  [2]int{1, 1},
		Attack: [2]int{100, 100},
	},
}

var WorldCharacterSpecs = []CharacterSpec{
	CharacterSpec{
		Name:  "Lion",
		Level: [2]int{1, 2},
		XP:    [2]float64{1, 1000},
		HP:    [2]int{1, 100},
		MaxHP: [2]int{100, 100},
	},
	CharacterSpec{
		Name:  "Liger",
		Level: [2]int{2, 1000},
		XP:    [2]float64{1, 1000},
		HP:    [2]int{1, 100},
		MaxHP: [2]int{100, 100},
	},
}

func (item *Item) GetSellPriceRange() string {
	switch item.Name {
	case WOODEN_SWORD:
		if item.Level == 1 { // attack 3
			return "60-63"
		} else if item.Level == 2 { // attack 6
			return "120-126"
		}
	case COPPER_SWORD:
		if item.Level == 1 { // attack 10
			return "200-210"
		} else if item.Level == 2 { // attack 20
			return "400-440"
		}
	}
	return "-1"
}

func (item *Item) GetUpgradePrice() int {
	switch item.Name {
	case WOODEN_SWORD:
		return 100
	case COPPER_SWORD:
		return 250
	}
	return -1
}
