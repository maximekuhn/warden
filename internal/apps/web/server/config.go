package server

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	MinecraftServers minecraftServers `yaml:"minecraftServers"`
}

type minecraftServers struct {
	Hostname       string                  `yaml:"hostname"`
	PortAllocation mcServersPortAllocation `yaml:"portAllocation"`
}

type mcServersPortAllocation struct {
	Strategy string   `yaml:"strategy"`
	Ports    []uint16 `yaml:"ports"`
}

func ParseConfigFromFile(filepath string) (*Config, error) {
	confFile, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := yaml.Unmarshal(confFile, &c); err != nil {
		return nil, err
	}
	return &c, c.Validate()
}

func (c *Config) Validate() error {
	if c.MinecraftServers.PortAllocation.Strategy != "pre-allocated" {
		return errors.New("for now, only portAllocation.strategy.'pre-allocated' is supported")
	}
	if len(c.MinecraftServers.PortAllocation.Ports) < 1 {
		return errors.New("portAllocation.strategy.'pre-allocated' strategy requires at least one port")
	}
	return nil
}
