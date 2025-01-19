package valueobjects

type MinecraftServerName struct {
	val string
}

func NewMinecraftServerName(val string) (MinecraftServerName, error) {
	// TODO: validation
	return MinecraftServerName{val: val}, nil
}
