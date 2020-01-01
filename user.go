package loud

// User represents an active user in the system.
type User interface {
	Reload()
	Save()
}

// UserSSHAuthentication for storing SSH auth.
type UserSSHAuthentication interface {
	SSHKeysEmpty() bool
	ValidateSSHKey(string) bool
	AddSSHKey(string)
}
