package transactions_test

import (
	"testing"

	"github.com/SebastianJ/elrond-sdk/transactions"
	"github.com/SebastianJ/elrond-sdk/utils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateShardForAddress(t *testing.T) {
	t.Parallel()

	numberOfShards := uint32(3)

	tests := []struct {
		address string
		shardID uint32
	}{
		{address: "erd1py99xhmfs4js6s6x0q6hp3c6zqwpq0fy4d3qmrwh7su68rudz2wq4e432r", shardID: 0},
		{address: "erd1trvgjh6c77j58wwqd5t7exdhxm2gsw6kwkf0xuzh0pza743jcjqqx0u87l", shardID: 0},
		{address: "erd1ffpa2ue77g50r4arz3rmqkxj3xykw4vgx7hyuxds9mc27ts97rtspfa6px", shardID: 1},
		{address: "erd1j6tzkx3dn2pu67tj83fuvgv4jnpsm3ehrw758lspkv9d03xjvpgs09trr4", shardID: 1},
		{address: "erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqplllst77y4l", shardID: transactions.MetachainShardID},
	}

	for _, test := range tests {
		addressBytes, err := utils.Bech32ToPublicKeyBytes(test.address)
		if err != nil {
			t.Error(err)
		}
		calculatedShardID := transactions.CalculateShardForAddress(addressBytes, numberOfShards)
		assert.Equal(t, test.shardID, calculatedShardID)
	}
}
