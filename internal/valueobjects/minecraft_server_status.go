package valueobjects

type MinecraftServerStatus string

const (
	MinecraftServerStatusRunning MinecraftServerStatus = "running"
	MinecraftServerStatusStopped MinecraftServerStatus = "stopped"
)
