package permissions

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Plan  Plan
	Roles map[uuid.UUID] /* server ID */ Role
}

func NewUser(
	id uuid.UUID,
	plan Plan,
	roles map[uuid.UUID] /* server ID */ Role,

) *User {
	return &User{
		ID:    id,
		Plan:  plan,
		Roles: roles,
	}
}
