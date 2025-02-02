package valueobjects

import "github.com/google/uuid"

type MinecraftServerID struct {
	val uuid.UUID
}

func NewMinecraftServerID(val uuid.UUID) (MinecraftServerID, error) {
	// TODO: validation
	msid := MinecraftServerID{}
	msid.val = val
	return msid, nil
}

func GenerateMinecraftServerID() MinecraftServerID {
	// no a correct value object, as it should validate the uuid
	// it will work for now (until a server id provider is created)
	return MinecraftServerID{
		val: uuid.New(),
	}
}

func (id MinecraftServerID) Value() uuid.UUID {
	return id.val
}
