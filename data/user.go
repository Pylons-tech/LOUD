package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
	MARKET
	SETTINGS
	DEVELOP
)

// User represents an active user in the system.
type User interface {
	SetGold(int)
	SetPylonAmount(int)
	SetItems([]Item)
	SetCharacters([]Character)
	SetDefaultItemIndex(idx int)
	SetDefaultCharacterIndex(idx int)
	SetLocation(UserLocation)
	SetLastTransaction(string)
	SetLatestBlockHeight(int64)
	GetLatestBlockHeight() int64
	InventoryItems() []Item
	InventoryItemIDByName(string) string
	InventoryIronSwords() []Item
	InventorySwords() []Item
	InventoryCharacters() []Character
	InventoryUpgradableItems() []Item
	InventorySellableItems() []Item
	GetLocation() UserLocation
	GetPrivKey() string
	GetDefaultItemIndex() int
	GetDefaultCharacterIndex() int
	GetDefaultCharacter() *Character
	GetGold() int
	GetPylonAmount() int
	GetUserName() string
	GetLastTransaction() string
	Reload()
	Save()
}
