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
	SetCharacters([]Item)
	SetDefaultItemIndex(idx int)
	SetDefaultCharacterIndex(idx int)
	SetLocation(UserLocation)
	SetLastTransaction(string)
	InventoryItems() []Item
	InventoryCharacters() []Item
	UpgradableItems() []Item
	GetLocation() UserLocation
	GetPrivKey() string
	GetDefaultItemIndex() int
	GetDefaultCharacterIndex() int
	GetDefaultCharacter() *Item
	GetGold() int
	GetPylonAmount() int
	GetUserName() string
	GetLastTransaction() string
	Reload()
	Save()
}
