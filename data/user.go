package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
	PYLCNTRL
	SETTINGS
	DEVELOP
	HELP
)

// User represents an active user in the system.
type User interface {
	SetAddress(string)
	SetGold(int)
	SetPylonAmount(int)
	SetItems([]Item)
	SetCharacters([]Character)
	SetActiveCharacterIndex(idx int)
	SetFightMonster(string)
	SetLocation(UserLocation)
	SetLastTransaction(string, string)
	SetLatestBlockHeight(int64)
	InventoryItems() []Item
	HasPreItemForAnItem(Item) bool
	InventoryItemIDByName(string) string
	InventoryAngelSwords() []Item
	InventoryIronSwords() []Item
	InventorySwords() []Item
	InventoryCharacters() []Character
	InventoryUpgradableItems() []Item
	InventorySellableItems() []Item
	GetLocation() UserLocation
	GetPrivKey() string
	GetActiveCharacterIndex() int
	GetActiveCharacter() *Character
	GetDeadCharacter() *Character
	GetFightWeapon() *Item
	GetItemByID(string) *Item
	GetAddress() string
	GetGold() int
	GetPylonAmount() int
	GetUserName() string
	GetLastTxHash() string
	GetLastTxMetaData() string
	GetLatestBlockHeight() int64
	Reload()
	Save()
}
