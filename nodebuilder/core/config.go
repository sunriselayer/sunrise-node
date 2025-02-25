package core

import (
	"fmt"
	"strconv"

	"github.com/sunriselayer/sunrise-da/libs/utils"
)

var MetricsEnabled bool

// Config combines all configuration fields for managing the relationship with a Core node.
type Config struct {
	IP       string
	RPCPort  string
	GRPCPort string
}

// DefaultConfig returns default configuration for managing the
// node's connection to a Celestia-Core endpoint.
func DefaultConfig() Config {
	return Config{
		IP:       "",
		RPCPort:  "26657",
		GRPCPort: "9090",
	}
}

// Validate performs basic validation of the config.
func (cfg *Config) Validate() error {
	if !cfg.IsEndpointConfigured() {
		return nil
	}

	ip, err := utils.ValidateAddr(cfg.IP)
	if err != nil {
		return err
	}
	cfg.IP = ip
	_, err = strconv.Atoi(cfg.RPCPort)
	if err != nil {
		return fmt.Errorf("nodebuilder/core: invalid rpc port: %s", err.Error())
	}
	_, err = strconv.Atoi(cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("nodebuilder/core: invalid grpc port: %s", err.Error())
	}
	return nil
}

// IsEndpointConfigured returns whether a core endpoint has been set
// on the config (true if set).
func (cfg *Config) IsEndpointConfigured() bool {
	return cfg.IP != ""
}
