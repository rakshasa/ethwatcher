package blockchain

import (
	"fmt"
	"math/big"

	"github.com/rakshasa/ethwatcher/utils"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Log struct {
	ethtypes.Log
	topics []Topic
}

func NewLog(log *ethtypes.Log) *Log {
	topics := make([]Topic, len(log.Topics))

	for idx, topic := range log.Topics {
		topics[idx] = Topic{topic}
	}

	return &Log{
		Log:    *log,
		topics: topics,
	}
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

func (l *Log) AddressAsHex() string {
	return fmt.Sprintf("0x%040x", new(big.Int).SetBytes(l.Address.Bytes()))
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

func (l *Log) Topics() []Topic {
	return l.topics
}

func (l *Log) TopicsAsBig() []big.Int {
	results := make([]big.Int, len(l.topics))

	for idx, topic := range l.topics {
		results[idx] = *topic.Big()
	}

	return results
}
