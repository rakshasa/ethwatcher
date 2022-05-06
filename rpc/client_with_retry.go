package rpc

import (
	"context"
	"time"

	"github.com/rakshasa/ethwatcher/blockchain"
)

type clientWithRetry struct {
	*client
	maxRetries uint
}

func (rpc *clientWithRetry) NewFilter(addresses []string, topics [][]string) (filterId string, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		filterId, err = rpc.client.NewFilter(addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) BlockByNumWithoutTx(ctx context.Context, blockNum uint64) (block *blockchain.Block, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		block, err = rpc.client.BlockByNumWithoutTx(ctx, blockNum)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) BlockNumber(ctx context.Context) (result uint64, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		result, err = rpc.client.BlockNumber(ctx)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) GetBlockByNum(num uint64) (rst *blockchain.Block, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.client.GetBlockByNum(num)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) GetFilterChanges(filterId string) (logs []blockchain.IReceiptLog, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		logs, err = rpc.client.GetFilterChanges(filterId)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) GetLogs(fromBlockNum, toBlockNum uint64, addresses []string, topics [][]string) (rst []blockchain.IReceiptLog, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.client.GetLogs(fromBlockNum, toBlockNum, addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) GetTransactionReceipt(txHash string) (rst blockchain.TransactionReceipt, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		rst, err = rpc.client.GetTransactionReceipt(txHash)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}
