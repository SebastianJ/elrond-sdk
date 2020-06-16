package transactions

import (
	"math"

	"github.com/ElrondNetwork/elrond-go/core"
)

// CalculateShardForAddress - calculates the shard for a given address
func CalculateShardForAddress(address []byte, numberOfShards uint32) uint32 {
	bytesNeed := int(numberOfShards/256) + 1
	startingIndex := 0
	if len(address) > bytesNeed {
		startingIndex = len(address) - bytesNeed
	}

	buffNeeded := address[startingIndex:]
	if core.IsSmartContractOnMetachain(buffNeeded, address) {
		return core.MetachainShardId
	}

	addr := uint32(0)
	for i := 0; i < len(buffNeeded); i++ {
		addr = addr<<8 + uint32(buffNeeded[i])
	}

	maskHigh, maskLow := calculateMasks(numberOfShards)

	shard := addr & maskHigh
	if shard > numberOfShards-1 {
		shard = addr & maskLow
	}

	return shard
}

func calculateMasks(numberOfShards uint32) (uint32, uint32) {
	n := math.Ceil(math.Log2(float64(numberOfShards)))
	return (1 << uint(n)) - 1, (1 << uint(n-1)) - 1
}
