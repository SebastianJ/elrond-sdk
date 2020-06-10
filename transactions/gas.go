package transactions

import (
	"strconv"

	"github.com/SebastianJ/elrond-sdk/config"
)

var (
	// DefaultGasParams - default gas params when gas params can't be used from economics.toml
	DefaultGasParams = GasParams{
		GasPrice:       200000000000,
		GasLimit:       50000,
		GasPerDataByte: 1500,
	}
)

// GasParams - represents gas parameters for a transaction
type GasParams struct {
	GasPrice       uint64
	GasLimit       uint64
	GasPerDataByte uint64
}

// ParseGasSettings - parse relevant gas settings from economics.toml
func ParseGasSettings(configPath string) (GasParams, error) {
	economicsConfig, err := config.LoadEconomicsConfig(configPath)
	if err != nil {
		return DefaultGasParams, err
	}

	gasPrice, err := strconv.ParseInt(economicsConfig.FeeSettings.MinGasPrice, 10, 64)
	if err != nil {
		return DefaultGasParams, err
	}

	gasLimit, err := strconv.ParseInt(economicsConfig.FeeSettings.MinGasLimit, 10, 64)
	if err != nil {
		return DefaultGasParams, err
	}

	gasPerDataByte, err := strconv.ParseInt(economicsConfig.FeeSettings.GasPerDataByte, 10, 64)
	if err != nil {
		return DefaultGasParams, err
	}

	gasParams := GasParams{
		GasPrice:       uint64(gasPrice),
		GasLimit:       uint64(gasLimit),
		GasPerDataByte: uint64(gasPerDataByte),
	}

	return gasParams, nil
}
