package valueobjects

import (
	"errors"
	"fmt"
	"net"
)

type MinecraftServerHostConfig struct {
	IPv4 net.IPAddr
	Port uint16
}

func NewMinecraftServerMetadata(
	ipv4 net.IPAddr,
	port uint16,
) (MinecraftServerHostConfig, error) {
	if ipv4.IP.DefaultMask() == nil {
		return MinecraftServerHostConfig{}, fmt.Errorf("%s is not a valid IPv4", ipv4.String())
	}
	if port < 5000 {
		return MinecraftServerHostConfig{}, errors.New("port must be > 5000")
	}
	return MinecraftServerHostConfig{
		IPv4: ipv4,
		Port: port,
	}, nil

}
