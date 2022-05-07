package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rakshasa/ethwatcher/blockchain"
)

type client struct {
	client *ethclient.Client
}

func (rpc *client) NewFilter(addresses []string, topics [][]string) (string, error) {
	// filterParams := ethrpc.FilterParams{
	// 	FromBlock: "latest",
	// 	ToBlock:   "latest",
	// 	Address:   addresses,
	// 	Topics:    topics,
	// }

	return "", fmt.Errorf("not implemented")

	// filterId, err := rpc.rpcImpl.EthNewFilter(filterParams)
	// if err != nil {
	// 	utils.Warnf("eth_newfilter: failed to create new filter: %v", err)
	// 	return "", err
	// }

	// return filterId, err
}

func (rpc *client) BlockByNumber(ctx context.Context, blockNumber uint64) (*blockchain.Block, error) {
	block, err := rpc.client.BlockByNumber(ctx, new(big.Int).SetUint64(blockNumber))
	if err != nil {
		return nil, fmt.Errorf("rpc request failed: %v", err)
	}

	return blockchain.NewBlock(block), nil
}

func (rpc *client) BlockNumber(ctx context.Context) (uint64, error) {
	return rpc.client.BlockNumber(ctx)
}

// BlockHash *common.Hash       used by eth_getLogs, return logs only from block with this hash
// FromBlock *big.Int           beginning of the queried range, nil means genesis block
// ToBlock   *big.Int           end of the range, nil means latest block
// Addresses []common.Address   restricts matches to events created by specific contracts
//
// The Topic list restricts matches to particular event topics. Each event has a list
// of topics. Topics matches a prefix of that list. An empty element slice matches any
// topic. Non-empty elements represent an alternative that matches any of the
// contained topics.
//
// Examples:
// {} or nil          matches any topic list
// {{A}}              matches topic A in first position
// {{}, {B}}          matches any topic in first position AND B in second position
// {{A}, {B}}         matches topic A in first position AND B in second position
// {{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
func (rpc *client) FilterLogs(ctx context.Context, fromBlock, toBlock uint64, addresses []string, topics [][]string) ([]blockchain.Log, error) {
	filterQuery := ethereum.FilterQuery{
		BlockHash: nil,
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
	}

	logs, err := rpc.client.FilterLogs(ctx, filterQuery)
	if err != nil {
		return nil, fmt.Errorf("rpc request failed: %v", err)
	}

	results := make([]blockchain.Log, len(logs))

	for idx, log := range logs {
		results[idx] = *blockchain.NewLog(&log)
	}

	return results, nil
}

func (rpc *client) GetFilterChanges(filterId string) ([]blockchain.Log, error) {
	return nil, fmt.Errorf("not implemented")

	// logs, err := rpc.rpcImpl.EthGetFilterChanges(filterId)
	// if err != nil {
	// 	utils.Warnf("eth_getfilterchanges: failed to retrieve filter changes: %v", err)
	// 	return nil, err
	// }

	// var result []*blockchain.Log
	// for i := 0; i < len(logs); i++ {
	// 	l := logs[i]
	// 	result = append(result, blockchain.ReceiptLog{Log: &l})

	// 	utils.Tracef("eth_getlogs: receipt log: %+v", l)
	// }

	// return result, err
}

// func (rpc *client) GetTransactionReceipt(txHash string) (blockchain.TransactionReceipt, error) {
// 	return nil, fmt.Errorf("not implemented")

// receipt, err := rpc.rpcImpl.EthGetTransactionReceipt(txHash)
// if err != nil {
// 	return nil, err
// }
// if receipt == nil {
// 	return nil, errors.New("nil receipt")
// }

// return &blockchain.EthereumTransactionReceipt{receipt}, err
// }
