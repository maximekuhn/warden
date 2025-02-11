package valueobjects

type MinecraftServerStatus string

const (
	MinecraftServerStatusRunning  MinecraftServerStatus = "running"
	MinecraftServerStatusStarting MinecraftServerStatus = "starting"
	MinecraftServerStatusStopped  MinecraftServerStatus = "stopped"
)
