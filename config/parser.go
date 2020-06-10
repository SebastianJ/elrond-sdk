package config

import (
	"github.com/ElrondNetwork/elrond-go/config"
	"github.com/ElrondNetwork/elrond-go/core"
)

// LoadEconomicsConfig - loads the economics.toml file to determine gas settings
func LoadEconomicsConfig(filepath string) (*config.EconomicsConfig, error) {
	cfg := &config.EconomicsConfig{}
	err := core.LoadTomlFile(cfg, filepath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
