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
		},
		Data: data,
	})

	t0a, ok := v.TopicAtIndexAsHex(0)
	if assert.True(ok) {
		assert.Equal("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", t0a)
	}

	t0ah, ok := v.TopicAtIndexAsAddressHex(0)
	assert.False(ok)
	assert.Empty(t0ah)

	t1ah, ok := v.TopicAtIndexAsAddressHex(1)
	if assert.True(ok) {
		assert.Equal("0x8032eaede5c55f744387ca53aaf0499abcd783e5", t1ah)
	}

	d0, ok := v.DataAtIndexAsUint64(0)
	if assert.True(ok) {
		assert.Equal(uint64(0x1b7f211b33c7c200), d0)
	}
}
