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

func (l *Log) TopicsAsBig() []big.Int {
	results := make([]big.Int, len(l.Topics))

	for idx, t := range l.Topics {
		results[idx] = *t.Big()
	}

	return results
}
