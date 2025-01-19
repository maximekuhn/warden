package permissions

// An action is associated to one or more Role(s) and represent what the associated
// Role(s) can perform within a specific minecraft server.
type Action string

const (
	ActionStartServer Action = "action.startServer"
	ActionStopServer  Action = "action.stopServer"
)
