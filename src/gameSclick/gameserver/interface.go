package gameserver

// Server defines the actions every game server must support
type Server interface {
	Install() error
	GetPath() string
	GetName() string
}
