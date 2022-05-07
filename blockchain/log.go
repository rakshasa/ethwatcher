package blockchain

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Log struct {
	ethtypes.Log
}

func NewLog(log *ethtypes.Log) *Log {
	return &Log{*log}
}
