package transactions_test

import (
	"testing"

	"github.com/SebastianJ/elrond-sdk/crypto"
	"github.com/SebastianJ/elrond-sdk/transactions"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCreateValidatorPayload(t *testing.T) {
	t.Parallel()

	blsPublicKeys := []crypto.Key{
		crypto.Key{PublicKeyString: "c4ac6bc7b77126a793d3239f8a13c24a93af191304ae423006780c4c6a357be71fabce8f239c8183fe0c3d7d51c89b1173c6c7bf272e0e044f0e7cbc46d613c698ca840c25a729c139cc4cbefd9b5f9745ccca62ab36d1b45b1d9952ffb28b8f"},
		crypto.Key{PublicKeyString: "5a8973ddd15907de5ca382c60fe3b3f66db1f18b41a1515199cd9b07606b75bb1ce8aa64df99e8a5466131cb772b0402ba9b0950b286a86a6d468d642620c0412ec75fbb513b208eee8bf536471a8df67957ca713acbe5be25fd577a18b26f14"},
		crypto.Key{PublicKeyString: "b827f28da78678d905f2a79512ae6724f5807710d079637a9907d89b26427b09c2dde46fae1bf6c0aa3749aac06328061cbf21efc29984912184d60e81dcd4fe9f839ae8fa4b6c2f7d6583d1ded7d40a1e2600f21feb02ea9254b2a2c7c60491"},
	}

	payload := transactions.GenerateCreateValidatorPayload(blsPublicKeys, "abc123")
	expectedPayload := "stake@3@c4ac6bc7b77126a793d3239f8a13c24a93af191304ae423006780c4c6a357be71fabce8f239c8183fe0c3d7d51c89b1173c6c7bf272e0e044f0e7cbc46d613c698ca840c25a729c139cc4cbefd9b5f9745ccca62ab36d1b45b1d9952ffb28b8f@abc123@5a8973ddd15907de5ca382c60fe3b3f66db1f18b41a1515199cd9b07606b75bb1ce8aa64df99e8a5466131cb772b0402ba9b0950b286a86a6d468d642620c0412ec75fbb513b208eee8bf536471a8df67957ca713acbe5be25fd577a18b26f14@abc123@b827f28da78678d905f2a79512ae6724f5807710d079637a9907d89b26427b09c2dde46fae1bf6c0aa3749aac06328061cbf21efc29984912184d60e81dcd4fe9f839ae8fa4b6c2f7d6583d1ded7d40a1e2600f21feb02ea9254b2a2c7c60491@abc123"

	assert.Equal(t, payload, expectedPayload)
}

func TestGenerateCreateValidatorPayloadUsingRandomBlsKeys(t *testing.T) {
	numBlsKeys := 3
	blsKeys, err := crypto.GenerateBlsKeys(numBlsKeys)
	if err != nil {
		t.Error(err)
	}
	payload := transactions.GenerateCreateValidatorPayload(blsKeys, "abc123")

	assert.Len(t, payload, 607)

	//t.Errorf("%s", payload)
}
