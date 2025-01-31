package services

import "net"

type MinecraftServerMetadataService interface {
	GetIPv4Addr() net.IPAddr
	ListAvailablePort() ([]uint16, error)
}
