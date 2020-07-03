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

// Friend is a struct to manage friend
type Friend struct {
	Name    string
	Address string
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

	// WoodenSword is string for wooden sword
	WoodenSword = "Wooden sword"
	// CopperSword is string for copper sword
	CopperSword = "Copper sword"
	// SilverSword is string for silver sword
	SilverSword = "Silver sword"
	// BronzeSword is string for bronze sword
	BronzeSword = "Bronze sword"
	// IronSword is string for iron sword
	IronSword = "Iron sword"
	// AngelSword is string for angel sword
	AngelSword = "Angel sword"

	// GoblinEar is string for goblin ear
	GoblinEar = "Goblin ear"
	// GoblinBoots is string for goblin boots
	GoblinBoots = "Goblin boots"
	// WolfFur is string for wolf fur
	WolfFur = "Wolf fur"
	// TrollSmellyBones is string for troll smelly bones
	TrollSmellyBones = "Troll smelly bones"
	// WolfTail is string for wolf tail
	WolfTail = "Wolf tail"
	// TrollToes is string for troll toes
	TrollToes = "Troll toes"
	// DropDragonIce is string for icy shards
	DropDragonIce = "Icy shards"
	// DropDragonFire is string for Fire scale
	DropDragonFire = "Fire scale"
	// DropDragonAcid is string for Poison claws
	DropDragonAcid = "Poison claws"
)

// IsSword is a helper function to check if an item is a kind of sword item
func (item Item) IsSword() bool {
	return strings.Contains(item.Name, "sword")
}

// ShopItems describes the items that are buyable at shop
var ShopItems = []Item{
	{
		ID:    "001",
		Name:  WoodenSword,
		Level: 1,
		Price: 100,
	},
	{
		ID:    "002",
		Name:  CopperSword,
		Level: 1,
		Price: 250,
	},
	{
		ID:       "003",
		Name:     SilverSword,
		Level:    1,
		Price:    50,
		PreItems: []string{GoblinEar},
	},
	{
		ID:       "004",
		Name:     BronzeSword,
		Level:    1,
		Price:    10,
		PreItems: []string{WolfTail},
	},
	{
		ID:       "005",
		Name:     IronSword,
		Level:    1,
		Price:    250,
		PreItems: []string{TrollToes},
	},
	{
		ID:       "006",
		Name:     AngelSword,
		Level:    1,
		Price:    20000,
		PreItems: []string{DropDragonFire, DropDragonIce, DropDragonAcid},
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
	{
		Name:   WoodenSword,
		Level:  [2]int{1, 1},
		Attack: [2]int{3, 3},
	},
	{
		Name:   WoodenSword,
		Level:  [2]int{2, 2},
		Attack: [2]int{6, 6},
	},
	{
		Name:   CopperSword,
		Level:  [2]int{1, 1},
		Attack: [2]int{10, 10},
	},
	{
		Name:   CopperSword,
		Level:  [2]int{2, 2},
		Attack: [2]int{20, 20},
	},
	{
		Name:   SilverSword,
		Level:  [2]int{1, 1},
		Attack: [2]int{30, 30},
	},
	{
		Name:   BronzeSword,
		Level:  [2]int{1, 1},
		Attack: [2]int{50, 50},
	},
	{
		Name:   IronSword,
		Level:  [2]int{1, 1},
		Attack: [2]int{100, 100},
	},
	{
		Name:  TrollToes,
		Level: [2]int{1, 1},
	},
	{
		Name:  WolfTail,
		Level: [2]int{1, 1},
	},
	{
		Name:  GoblinEar,
		Level: [2]int{1, 1},
	},
	{
		Name:   AngelSword,
		Level:  [2]int{1, 1},
		Attack: [2]int{1000, 1000},
	},
	{
		Name:  DropDragonFire,
		Level: [2]int{1, 1},
	},
	{
		Name:  DropDragonIce,
		Level: [2]int{1, 1},
	},
	{
		Name:  DropDragonAcid,
		Level: [2]int{1, 1},
	},
}

// WorldCharacterSpecs characters that are buyable from pylons central by trading
var WorldCharacterSpecs = []CharacterSpec{
	{
		Name:  "LionBaby",
		Level: [2]int{1, 2},
		XP:    [2]float64{1, 1000000},
	},
	{
		Special: FireSpecial,
		Name:    "FireBaby",
		Level:   [2]int{1, 1000},
		XP:      [2]float64{1, 1000000},
	},
	{
		Special: IceSpecial,
		Name:    "IceBaby",
		Level:   [2]int{1, 1000},
		XP:      [2]float64{1, 1000000},
	},
	{
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
	case WoodenSword:
		return 100
	case CopperSword:
		return 250
	}
	return -1
}
