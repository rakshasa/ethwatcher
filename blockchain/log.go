package blockchain

import (
	"math/big"

	"github.com/rakshasa/ethwatcher/utils"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Log struct {
	ethtypes.Log
}

func NewLog(log *ethtypes.Log) *Log {
	return &Log{*log}
}

func NewTopicFromHex(str string) (ethcommon.Hash, error) {
	value, err := utils.NewBigFromHex(str, 8*ethcommon.HashLength)
	if err != nil {
		return ethcommon.Hash{}, err
	}

	return ethcommon.BigToHash(value), nil
}

func (l *Log) AddressAsBig() *big.Int {
	return new(big.Int).SetBytes(l.Address.Bytes())
}

// TODO: Change to use a custom type for data.

func (l *Log) DataLen() int {
	return len(l.Data) / 32
}

func (l *Log) DataAtIndexAsBig(idx int) (*big.Int, bool) {
	if idx < 0 || idx >= len(l.Data)/32 {
		return nil, false
	}

	return new(big.Int).SetBytes(l.Data[idx*32 : (idx+1)*32]), true
}

func (l *Log) DataAtIndexAsUint64(idx int) (uint64, bool) {
	value, ok := l.DataAtIndexAsBig(idx)
	if !ok {
		return uint64(0), false
	}
	if !value.IsUint64() {
		return uint64(0), false
	}

	return value.Uint64(), true
}

func (l *Log) TopicsAsBig() []big.Int {
	results := make([]big.Int, len(l.Topics))

	for idx, t := range l.Topics {
		results[idx] = *t.Big()
	}

	return results
}

func (l *Log) TopicAtIndexAsAddressHex(idx int) (string, bool) {
	if idx < 0 || idx >= len(l.Topics) {
		return "", false
	}

	value := l.Topics[idx].Big()

	if value.BitLen() > 8*ethcommon.AddressLength {
		return "", false
	}

	return "0x" + value.Text(16), true
}

func (l *Log) TopicAtIndexAsHex(idx int) (string, bool) {
	if idx < 0 || idx >= len(l.Topics) {
		return "", false
	}

	return "0x" + l.Topics[idx].Big().Text(16), true
}

func (l *Log) TopicAtIndexAsUint64(idx int) (uint64, bool) {
	if idx < 0 || idx >= len(l.Topics) {
		return uint64(0), false
	}

	value := l.Topics[idx].Big()

	if value.IsUint64() {
		return uint64(0), false
	}

	return value.Uint64(), true
}
