package rpc

import (
	"github.com/rakshasa/ethwatcher/blockchain"
)

type IBlockChainRPC interface {
	NewFilter(addresses []string, topics [][]string) (string, error)

	GetBlockByNum(uint64) (blockchain.Block, error)
	GetBlockByNumWithoutTx(uint64) (blockchain.Block, error)
	GetCurrentBlockNum() (uint64, error)
	GetFilterChanges(filterId string) ([]blockchain.IReceiptLog, error)
	GetLogs(from, to uint64, address []string, topics [][]string) ([]blockchain.IReceiptLog, error)
	GetTransactionReceipt(txHash string) (blockchain.TransactionReceipt, error)
}
