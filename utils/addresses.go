package utils

import (
	"encoding/hex"
	"strings"

	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
)

// IdentifyAddressShard - identifies what shard an address belongs to
// This isn't working right now - based on old previous shard detection
func IdentifyAddressShard(address string) (int, error) {
	if strings.HasPrefix(address, "erd") {
		pubAddress, err := Bech32ToPublicKey(address)
		if err != nil {
			return -1, err
		}
		address = pubAddress
	}

	lastAddressCharacter := strings.ToLower(string(address[len(address)-1]))

	var shard int

	switch lastAddressCharacter {
	case "0", "8":
		shard = 0
	case "1", "5", "9", "d":
		shard = 1
	case "2", "6", "a", "e":
		shard = 2
	case "3", "7", "b", "f":
		shard = 3
	case "4", "c":
		shard = 4
	}

	return shard, nil
}

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
