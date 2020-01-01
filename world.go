package loud

// World represents a gameplay world. It should keep track of the map,
// entities in the map, and players.
type World interface {
	GetUser(string) User
	Close()
}
