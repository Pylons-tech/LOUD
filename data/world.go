package loud

// World represents a gameplay world.
type World interface {
	GetUser(string) User
	Close()
}

var SomethingWentWrongMsg string = ""
