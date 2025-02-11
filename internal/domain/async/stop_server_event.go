package async

import "github.com/maximekuhn/warden/internal/domain/valueobjects"

type StopServerEvent struct {
	ServerID valueobjects.MinecraftServerID
}
