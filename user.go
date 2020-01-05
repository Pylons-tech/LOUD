package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
)

// User represents an active user in the system.
type User interface {
	GetLocation() UserLocation
	Reload()
	Save()
}
