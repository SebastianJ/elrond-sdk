package crypto_test

import (
	"testing"

	"github.com/SebastianJ/elrond-sdk/crypto"
	"github.com/stretchr/testify/assert"
)

func TestBlsKeysGeneration(t *testing.T) {
	t.Parallel()

	numBlsKeys := 10
	blsKeys, err := crypto.GenerateBlsKeys(numBlsKeys)

	assert.Nil(t, err)
	assert.Len(t, blsKeys, numBlsKeys)
	assert.IsType(t, crypto.Key{}, blsKeys[0])
}

func TestBlsKeyGeneration(t *testing.T) {
	t.Parallel()

	blsKey, err := crypto.GenerateBlsKey()

	assert.Nil(t, err)
	assert.IsType(t, crypto.Key{}, blsKey)
}
