package rpc

import (
	"time"

	"github.com/onrik/ethrpc"
	"github.com/rakshasa/ethwatcher/blockchain"
)

type EthBlockChainRPCWithRetry struct {
	*EthBlockChainRPC
	maxRetryTimes int
}

func NewEthRPCWithRetry(api string, maxRetryCount int, options ...func(rpc *ethrpc.EthRPC)) *EthBlockChainRPCWithRetry {
	rpc := NewEthRPC(api, options...)

	return &EthBlockChainRPCWithRetry{rpc, maxRetryCount}
}

func (rpc EthBlockChainRPCWithRetry) NewFilter(addresses []string, topics []string) (filterId string, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		filterId, err = rpc.EthBlockChainRPC.NewFilter(addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetBlockByNum(num uint64) (rst blockchain.Block, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.EthBlockChainRPC.GetBlockByNum(num)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetBlockByNumWithoutTx(num uint64) (rst blockchain.Block, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.EthBlockChainRPC.GetBlockByNumWithoutTx(num)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetCurrentBlockNum() (rst uint64, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.EthBlockChainRPC.GetCurrentBlockNum()
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetFilterChanges(filterId string) (logs []blockchain.IReceiptLog, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		logs, err = rpc.EthBlockChainRPC.GetFilterChanges(filterId)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetLogs(
	fromBlockNum, toBlockNum uint64,
	addresses []string,
	topics []string,
) (rst []blockchain.IReceiptLog, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.EthBlockChainRPC.GetLogs(fromBlockNum, toBlockNum, addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc EthBlockChainRPCWithRetry) GetTransactionReceipt(txHash string) (rst blockchain.TransactionReceipt, err error) {
	for i := 0; i <= rpc.maxRetryTimes; i++ {
		rst, err = rpc.EthBlockChainRPC.GetTransactionReceipt(txHash)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}
