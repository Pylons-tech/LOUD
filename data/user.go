package loud

// UserLocation is a struct to manage user's location
type UserLocation int

// page tabs
const (
	// Home is home
	Home UserLocation = iota
	// Forest is a place to fight monsters
	Forest
	// Shop is a place to do shopping
	Shop
	// PylonsCentral is a place to do trading and buying from pylons
	PylonsCentral
	// Friends is a place to manage friends
	Friends
	// Settings is a place to control game settings
	Settings
	// Develop is a place to run development functions
	Develop
	// Help is a place to get help text
	Help
)

// User represents an active user in the system.
type User interface {
	SetAddress(string)
	SetGold(int)
	SetPylonAmount(int)
	SetLockedGold(int)
	SetLockedPylonAmount(int)
	SetItems([]Item)
	SetCharacters([]Character)
	SetActiveCharacterIndex(idx int)
	SetFightMonster(string)
	SetLocation(UserLocation)
	SetLastTransaction(string, string)
	SetLatestBlockHeight(int64)
	FixLoadedData()
	InventoryItems() []Item
	HasPreItemForAnItem(Item) bool
	InventoryItemIDByName(string) string
	InventoryAngelSwords() []Item
	InventoryIronSwords() []Item
	InventorySwords() []Item
	InventoryCharacters() []Character
	InventoryUpgradableItems() []Item
	InventorySellableItems() []Item
	Friends() []Friend
	SetFriends([]Friend)
	GetMatchedItems(ItemSpec) []Item
	GetMatchedCharacters(CharacterSpec) []Character
	GetLocation() UserLocation
	GetPrivKey() string
	GetActiveCharacterIndex() int
	GetActiveCharacter() *Character
	GetDeadCharacter() *Character
	GetTargetMonster() string
	GetFightWeapon() *Item
	GetItemByID(string) *Item
	GetAddress() string
	GetGold() int
	GetLockedGold() int
	GetPylonAmount() int
	GetLockedPylonAmount() int
	GetUserName() string
	GetLastTxHash() string
	GetLastTxMetaData() string
	GetLatestBlockHeight() int64
	Reload()
	Save()
}
