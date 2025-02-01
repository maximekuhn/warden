package async

import "github.com/maximekuhn/warden/internal/domain/valueobjects"

type ServerStartedEvent struct {
	ServerID valueobjects.MinecraftServerID
}
