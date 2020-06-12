package transactions

import (
	"fmt"
	"strings"

	"github.com/SebastianJ/elrond-sdk/crypto"
)

var (
	defaultPrivateKeyPlaceholder = "abc123"
)

// GenerateCreateValidatorPayload - generate the payload used to create a validator
func GenerateCreateValidatorPayload(blsKeys []crypto.Key, privateKeyPlaceholder string) string {
	if privateKeyPlaceholder == "" {
		privateKeyPlaceholder = defaultPrivateKeyPlaceholder
	}

	keyCount := len(blsKeys)
	countPrefix := fmt.Sprintf("%04d", keyCount)

	var payload strings.Builder
	payload.WriteString("stake@")
	payload.WriteString(fmt.Sprintf("%s@", countPrefix))

	for index, blsKey := range blsKeys {
		payload.WriteString(fmt.Sprintf("%s@%s", blsKey.PublicKeyString, privateKeyPlaceholder))
		if index != keyCount-1 {
			payload.WriteString("@")
		}
	}

	return payload.String()
}
