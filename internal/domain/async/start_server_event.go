package async

import "github.com/maximekuhn/warden/internal/domain/valueobjects"

type StartServerEvent struct {
	ServerID valueobjects.MinecraftServerID
}
