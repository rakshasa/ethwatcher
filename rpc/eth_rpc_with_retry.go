package rpc

import (
	"context"
	"time"

	"github.com/onrik/ethrpc"
	"github.com/rakshasa/ethwatcher/blockchain"
)

type EthBlockChainRPCWithRetry struct {
	*Client
	maxRetries uint
}

// TODO: Add retry interval.
func NewEthRPCWithRetry(api string, maxRetries uint, options ...func(rpc *ethrpc.EthRPC)) *EthBlockChainRPCWithRetry {
	rpc, _ := Dial(api)

	return &EthBlockChainRPCWithRetry{
		rpc,
		maxRetries,
	}
}

func (rpc EthBlockChainRPCWithRetry) NewFilter(addresses []string, topics [][]string) (filterId string, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		filterId, err = rpc.Client.NewFilter(addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetBlockByNum(num uint64) (rst blockchain.Block, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.Client.GetBlockByNum(num)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetBlockByNumWithoutTx(num uint64) (rst blockchain.Block, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.Client.GetBlockByNumWithoutTx(num)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) BlockNumber(ctx context.Context) (result uint64, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		result, err = rpc.Client.BlockNumber(ctx)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetFilterChanges(filterId string) (logs []blockchain.IReceiptLog, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		logs, err = rpc.Client.GetFilterChanges(filterId)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetLogs(fromBlockNum, toBlockNum uint64, addresses []string, topics [][]string) (rst []blockchain.IReceiptLog, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.Client.GetLogs(fromBlockNum, toBlockNum, addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetTransactionReceipt(txHash string) (rst blockchain.TransactionReceipt, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.Client.GetTransactionReceipt(txHash)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}
