package loud

const (
	// TextTigerChr it's a constant to use on applications
	TextTigerChr string = "Tiger"
)

// ShopCharacters are characters that are buyable at blacksmith
var ShopCharacters = []Character{
	Character{
		ID:    "001",
		Name:  TextTigerChr,
		Level: 1,
		XP:    1,
		Price: 1, // in pylon
	},
}
