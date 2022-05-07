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

func (rpc *clientWithRetry) BlockByNumber(ctx context.Context, blockNumber uint64) (block *blockchain.Block, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		block, err = rpc.client.BlockByNumber(ctx, blockNumber)
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

func (rpc *clientWithRetry) FilterLogs(ctx context.Context, fromBlock, toBlock uint64, addresses []string, topics [][]string) (logs []blockchain.Log, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		logs, err = rpc.client.FilterLogs(ctx, fromBlock, toBlock, addresses, topics)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}

func (rpc *clientWithRetry) GetFilterChanges(filterId string) (logs []blockchain.Log, err error) {
	for i := uint(0); i <= rpc.maxRetries; i++ {
		logs, err = rpc.client.GetFilterChanges(filterId)
		if err == nil {
			break
		}

		time.Sleep(time.Duration(500*(i+1)) * time.Millisecond)
	}

	return
}
