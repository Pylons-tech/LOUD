package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
	SETTINGS
)

// User represents an active user in the system.
type User interface {
	SetGold(int)
	SetItems([]Item)
	GetGold() int
	GetUserName() string
	InventoryItems() []Item
	UpgradableItems() []Item
	GetLocation() UserLocation
	SetLocation(UserLocation)
	GetLastTransaction() string
	SetLastTransaction(string)
	Reload()
	Save()
}
