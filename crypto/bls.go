package crypto

import (
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
)

var (
	blsPublicKeyLength = 96
)

// GenerateBlsKeys - generate muultiple BLS keys at once
func GenerateBlsKeys(count int) ([]Key, error) {
	keys := []Key{}

	for i := 0; i < count; i++ {
		key, err := GenerateBlsKey()
		if err != nil {
			return keys, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// GenerateBlsKey - generates a new BLS key
func GenerateBlsKey() (Key, error) {
	key, err := GenerateKey(BLS12)
	if err != nil {
		return Key{}, err
	}

	converter, err := pubkeyConverter.NewHexPubkeyConverter(blsPublicKeyLength)
	if err != nil {
		return Key{}, err
	}

	key.Converter = converter
	key.PublicKeyString = key.Converter.Encode(key.PublicKeyBytes)

	return key, nil
}
