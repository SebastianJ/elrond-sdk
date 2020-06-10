package utils

import (
	"encoding/hex"

	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
)

// PublicKeyToBech32 - converts a public key to bech32 format
func PublicKeyToBech32(key string) (string, error) {
	converter, err := pubkeyConverter.NewBech32PubkeyConverter(32)
	if err != nil {
		return "", err
	}

	decoded, err := hex.DecodeString(key)
	if err != nil {
		return "", err
	}

	bech32 := converter.Encode(decoded)

	return bech32, nil
}

// Bech32ToPublicKey - converts a bech32 address to pub key format
func Bech32ToPublicKey(bech32 string) (string, error) {
	converter, err := pubkeyConverter.NewBech32PubkeyConverter(32)
	if err != nil {
		return "", err
	}

	decoded, err := converter.Decode(bech32)
	if err != nil {
		return "", err
	}

	pubKey := hex.EncodeToString(decoded)

	return pubKey, nil
}
