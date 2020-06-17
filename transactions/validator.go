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
	countPrefix := fmt.Sprintf("%x", keyCount)

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

// GenerateUnstakingPayload - Amount: 0 - Gas limit: 6000000
func GenerateUnstakingPayload(blsKey crypto.Key) string {
	return generateStakingPayload("unStake", blsKey)
}

// GenerateUnbondingPayload - Amount: 0 - Gas limit: 6000000
func GenerateUnbondingPayload(blsKey crypto.Key) string {
	return generateStakingPayload("unBond", blsKey)
}

// GenerateUnjailPayload - Amount: 2500 - Gas limit: 6000000
func GenerateUnjailPayload(blsKey crypto.Key) string {
	return generateStakingPayload("unJail", blsKey)
}

// GenerateChangeRewardAddressPayload - Amount: 0 - Gas limit: 6000000
func GenerateChangeRewardAddressPayload(hexPublicKey string) string {
	return fmt.Sprintf("changeRewardAddress@%s", hexPublicKey)
}

// GenerateClaimPayload - Amount: 0 - Gas limit: 6000000
func GenerateClaimPayload() string {
	return "claim"
}

func generateStakingPayload(command string, blsKey crypto.Key) string {
	return fmt.Sprintf("%s@%s", command, blsKey.PublicKeyString)
}
