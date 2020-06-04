package loud

import (
	"fmt"
	"strings"
)

// Item is a struct to manage game item attributes
type Item struct {
	ID         string `json:""`
	Name       string `json:""`
	Level      int
	Attack     int
	Price      int
	Value      int
	PreItems   []string
	LastUpdate int64
}

// ItemSpec is a struct to manage game item spec on trading
type ItemSpec struct {
	Name   string `json:""`
	Level  [2]int
	Attack [2]int
	Price  int
}

// Character is a struct to manage game character attributes
type Character struct {
	ID                string `json:""`
	Name              string `json:""`
	Level             int
	Price             int
	XP                float64
	GiantKill         int
	Special           int
	SpecialDragonKill int
	UndeadDragonKill  int
	LastUpdate        int64
}

// CharacterSpec is a struct to manage game character spec on trading
type CharacterSpec struct {
	Special int
	Name    string `json:""`
	Level   [2]int
	Price   int
	XP      [2]float64
}

const (
	// NoSpecial means character has no special
	NoSpecial = 0
	// FireSpecial means character has fire special
	FireSpecial = 1
	// IceSpecial means character has ice special
	IceSpecial = 2
	// AcidSpecial means character has acid special
	AcidSpecial = 3
)

const (
	// TextRabbit is constant for text Rabbit
	TextRabbit string = "Rabbit"
	// TextGoblin is constant for text Goblin
	TextGoblin = "Goblin"
	// TextWolf is constant for text Wolf
	TextWolf = "Wolf"
	// TextTroll is constant for text Troll
	TextTroll = "Troll"
	// TextGiant is constant for text Giant
	TextGiant = "Giant"
	// TextDragonFire is constant for text Fire Dragon
	TextDragonFire = "Fire dragon"
	// TextDragonIce is constant for text Ice Dragon
	TextDragonIce = "Ice dragon"
	// TextDragonAcid is constant for text Acid Dragon
	TextDragonAcid = "Acid dragon"
	// TextDragonUndead is constant for text Undead Dragon
	TextDragonUndead = "Undead dragon"

	WOODEN_SWORD = "Wooden sword"
	COPPER_SWORD = "Copper sword"
	SILVER_SWORD = "Silver sword"
	BRONZE_SWORD = "Bronze sword"
	IRON_SWORD   = "Iron sword"
	ANGEL_SWORD  = "Angel sword"

	GOBLIN_EAR         = "Goblin ear"
	GOBLIN_BOOTS       = "Goblin boots"
	WOLF_FUR           = "Wolf fur"
	TROLL_SMELLY_BONES = "Troll smelly bones"
	WOLF_TAIL          = "Wolf tail"
	TROLL_TOES         = "Troll toes"
	DROP_DRAGONICE     = "Icy shards"
	DROP_DRAGONFIRE    = "Fire scale"
	DROP_DRAGONACID    = "Poison claws"
)

// IsSword is a helper function to check if an item is a kind of sword item
func (item Item) IsSword() bool {
	return strings.Contains(item.Name, "sword")
}

// ShopItems describes the items that are buyable at shop
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
		ID:       "003",
		Name:     SILVER_SWORD,
		Level:    1,
		Price:    50,
		PreItems: []string{GOBLIN_EAR},
	},
	Item{
		ID:       "004",
		Name:     BRONZE_SWORD,
		Level:    1,
		Price:    10,
		PreItems: []string{WOLF_TAIL},
	},
	Item{
		ID:       "005",
		Name:     IRON_SWORD,
		Level:    1,
		Price:    250,
		PreItems: []string{TROLL_TOES},
	},
	Item{
		ID:       "006",
		Name:     ANGEL_SWORD,
		Level:    1,
		Price:    20000,
		PreItems: []string{DROP_DRAGONFIRE, DROP_DRAGONICE, DROP_DRAGONACID},
	},
}

// PreItemStr returns text that are required to make an item
func (item Item) PreItemStr() string {
	switch len(item.PreItems) {
	case 1:
		return fmt.Sprintf("\"%s\"", item.PreItems[0])
	case 3: // angel sword
		return fmt.Sprintf("\"%s\"", Localize("drops of 3 special dragons"))
	default:
		return ""
	}
}

// WorldItemSpecs describes the items that are buyable by trading
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
		Name:   ANGEL_SWORD,
		Level:  [2]int{1, 1},
		Attack: [2]int{1000, 1000},
	},
	ItemSpec{
		Name:  DROP_DRAGONFIRE,
		Level: [2]int{1, 1},
	},
	ItemSpec{
		Name:  DROP_DRAGONICE,
		Level: [2]int{1, 1},
	},
	ItemSpec{
		Name:  DROP_DRAGONACID,
		Level: [2]int{1, 1},
	},
}

// WorldCharacterSpecs characters that are buyable from pylons central by trading
var WorldCharacterSpecs = []CharacterSpec{
	CharacterSpec{
		Name:  "LionBaby",
		Level: [2]int{1, 2},
		XP:    [2]float64{1, 1000000},
	},
	CharacterSpec{
		Special: FireSpecial,
		Name:    "FireBaby",
		Level:   [2]int{1, 1000},
		XP:      [2]float64{1, 1000000},
	},
	CharacterSpec{
		Special: IceSpecial,
		Name:    "IceBaby",
		Level:   [2]int{1, 1000},
		XP:      [2]float64{1, 1000000},
	},
	CharacterSpec{
		Special: AcidSpecial,
		Name:    "AcidBaby",
		Level:   [2]int{1, 1000},
		XP:      [2]float64{1, 1000000},
	},
}

// GetSellPriceRange calculates sell price range based on item's value
func (item *Item) GetSellPriceRange() string {
	minPrice := item.Value * 8 / 10
	maxPrice := minPrice + 20
	return fmt.Sprintf("%d-%d", minPrice, maxPrice)
}

// GetUpgradePrice calculates the upgrade price based on item type
func (item *Item) GetUpgradePrice() int {
	switch item.Name {
	case WOODEN_SWORD:
		return 100
	case COPPER_SWORD:
		return 250
	}
	return -1
}
