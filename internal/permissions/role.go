package permissions

// Each Role is associated to a set of actions within a minecraft server.
// A user can have 1..* Role(s), but only 1 for a specific minecraft server.
//
// For instance, jeff can be 'viewer' in server-1, and admin in server-2
type Role string

const (
	RoleViewer Role = "role.viewer"
	RoleAdmin  Role = "role.admin"
)
