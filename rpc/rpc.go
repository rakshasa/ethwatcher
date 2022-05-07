package rpc

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rakshasa/ethwatcher/blockchain"
)

type Client interface {
	NewFilter(addresses []string, topics [][]string) (string, error)

	BlockNumber(ctx context.Context) (uint64, error)
	BlockByNumber(ctx context.Context, blockNumber uint64) (*blockchain.Block, error)
	FilterLogs(ctx context.Context, fromBlock, toBlock uint64, addresses []string, topics [][]string) ([]blockchain.Log, error)

	GetFilterChanges(filterId string) ([]blockchain.Log, error)
}

func Dial(rawurl string) (Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (Client, error) {
	c, err := ethclient.DialContext(ctx, rawurl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ethereum api server: %v", err)
	}

	return &client{
		client: c,
	}, nil
}

// TODO: Add retry interval.
func DialWithRetry(rawurl string, maxRetries uint) (Client, error) {
	return DialContextWithRetry(context.Background(), rawurl, maxRetries)
}

func DialContextWithRetry(ctx context.Context, rawurl string, maxRetries uint) (Client, error) {
	c, err := ethclient.DialContext(ctx, rawurl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ethereum api server: %v", err)
	}

	return &clientWithRetry{
		client: &client{
			c,
		},
		maxRetries: maxRetries,
	}, nil
}
