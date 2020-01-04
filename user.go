package loud

// User represents an active user in the system.
type User interface {
	Reload()
	Save()
}
