package rpc

import (
	"context"

	"github.com/rakshasa/ethwatcher/blockchain"
)

type IBlockChainRPC interface {
	NewFilter(addresses []string, topics [][]string) (string, error)

	BlockNumber(ctx context.Context) (uint64, error)

	GetBlockByNum(uint64) (blockchain.Block, error)
	GetBlockByNumWithoutTx(uint64) (blockchain.Block, error)
	GetFilterChanges(filterId string) ([]blockchain.IReceiptLog, error)
	GetLogs(from, to uint64, address []string, topics [][]string) ([]blockchain.IReceiptLog, error)
	GetTransactionReceipt(txHash string) (blockchain.TransactionReceipt, error)
}
