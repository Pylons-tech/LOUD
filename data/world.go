package loud

// World represents a gameplay world.
type World interface {
	GetUser(string) User
	Close()
}

// SomethingWentWrongMsg is a global variable to report something went wrong message to a user
var SomethingWentWrongMsg string = ""
