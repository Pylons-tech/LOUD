package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
)

// User represents an active user in the system.
type User interface {
	AddGold(amount int)
	GetGold() int
	GetUserName() string
	InventoryItems() []Item
	GetLocation() UserLocation
	SetLocation(UserLocation)
	GetLastTransaction() string
	SetLastTransaction(string)
	Reload()
	Save()
}
