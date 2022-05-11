package blockchain

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func TestLog(t *testing.T) {
	assert := assert.New(t)

	data, err := hex.DecodeString("0000000000000000000000000000000000000000000000001b7f211b33c7c200")
	if !assert.NoError(err) {
		return
	}

	v := NewLog(&ethtypes.Log{
		Topics: []ethcommon.Hash{
			ethcommon.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			ethcommon.HexToHash("0x0000000000000000000000008032eaede5c55f744387ca53aaf0499abcd783e5"),
			ethcommon.HexToHash("0x000000000000000000000000000000000000000000000000ffffffffffffffff"),
		},
		Data: data,
	})

	if !assert.Len(v.Topics(), 3) {
		return
	}

	assert.Equal("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", v.Topics()[0].Hex())

	t0, ok := v.Topics()[0].AddressHex()
	assert.False(ok)
	assert.Empty(t0)

	t1, ok := v.Topics()[1].AddressHex()
	assert.True(ok)
	assert.Equal("0x8032eaede5c55f744387ca53aaf0499abcd783e5", t1)

	t2ah, ok := v.Topics()[2].AddressHex()
	assert.True(ok)
	assert.Equal("0x000000000000000000000000ffffffffffffffff", t2ah)

	t2u32, ok := v.Topics()[2].Uint32()
	assert.False(ok)
	assert.Equal(uint32(0), t2u32)

	t2u64, ok := v.Topics()[2].Uint64()
	assert.True(ok)
	assert.Equal(uint64(0xffffffffffffffff), t2u64)

	d0, ok := v.DataAtIndexAsUint64(0)
	assert.True(ok)
	assert.Equal(uint64(0x1b7f211b33c7c200), d0)
}
