package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

func TestLogTopics(t *testing.T) {
	assert := assert.New(t)

	v := NewLog(&ethtypes.Log{
		Topics: []ethcommon.Hash{
			ethcommon.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			ethcommon.HexToHash("0x0000000000000000000000008032eaede5c55f744387ca53aaf0499abcd783e5"),
		},
	})

	t0ah, ok := v.TopicAtIndexAsAddressHex(0)
	assert.False(ok)
	assert.Empty(t0ah)

	t1ah, ok := v.TopicAtIndexAsAddressHex(1)
	if assert.True(ok) {
		assert.Equal("0x8032eaede5c55f744387ca53aaf0499abcd783e5", t1ah)
	}
}
