package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
	PYLCNTRL
	SETTINGS
	DEVELOP
)

// User represents an active user in the system.
type User interface {
	SetAddress(string)
	SetGold(int)
	SetPylonAmount(int)
	SetItems([]Item)
	SetCharacters([]Character)
	SetActiveWeaponIndex(idx int)
	SetActiveCharacterIndex(idx int)
	SetLocation(UserLocation)
	SetLastTransaction(string)
	SetLatestBlockHeight(int64)
	InventoryItems() []Item
	HasPreItemForAnItem(Item) bool
	InventoryItemIDByName(string) string
	InventoryIronSwords() []Item
	InventorySwords() []Item
	InventoryCharacters() []Character
	InventoryUpgradableItems() []Item
	InventorySellableItems() []Item
	GetLocation() UserLocation
	GetPrivKey() string
	GetActiveWeaponIndex() int
	GetActiveCharacterIndex() int
	GetActiveCharacter() *Character
	GetActiveWeapon() *Item
	GetAddress() string
	GetGold() int
	GetPylonAmount() int
	GetUserName() string
	GetLastTransaction() string
	GetLatestBlockHeight() int64
	Reload()
	Save()
}
