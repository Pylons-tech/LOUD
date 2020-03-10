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
	GetGold() int
	GetPylonAmount() int
	GetUserName() string
	InventoryItems() []Item
	InventoryCharacters() []Item
	UpgradableItems() []Item
	GetLocation() UserLocation
	GetPrivKey() string
	SetLocation(UserLocation)
	GetLastTransaction() string
	SetLastTransaction(string)
	Reload()
	Save()
}
