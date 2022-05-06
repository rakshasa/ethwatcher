package rpc

import (
	"context"
	"fmt"

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

func (rpc *client) GetBlockByNum(num uint64) (blockchain.Block, error) {
	return nil, fmt.Errorf("not implemented")

	// b, err := rpc.rpcImpl.EthGetBlockByNumber(int(num), true)
	// if err != nil {
	// 	return nil, err
	// }
	// if b == nil {
	// 	return nil, errors.New("nil block")
	// }

	// return &blockchain.EthereumBlock{b}, err
}

func (rpc *client) GetBlockByNumWithoutTx(num uint64) (blockchain.Block, error) {
	return nil, fmt.Errorf("not implemented")

	// b, err := rpc.rpcImpl.EthGetBlockByNumber(int(num), false)
	// if err != nil {
	// 	return nil, err
	// }
	// if b == nil {
	// 	return nil, errors.New("nil block")
	// }

	// return &blockchain.EthereumBlock{b}, err
}

func (rpc *client) BlockNumber(ctx context.Context) (uint64, error) {
	return rpc.client.BlockNumber(ctx)
}

func (rpc *client) GetFilterChanges(filterId string) ([]blockchain.IReceiptLog, error) {
	return nil, fmt.Errorf("not implemented")

	// logs, err := rpc.rpcImpl.EthGetFilterChanges(filterId)
	// if err != nil {
	// 	utils.Warnf("eth_getfilterchanges: failed to retrieve filter changes: %v", err)
	// 	return nil, err
	// }

	// var result []blockchain.IReceiptLog
	// for i := 0; i < len(logs); i++ {
	// 	l := logs[i]
	// 	result = append(result, blockchain.ReceiptLog{Log: &l})

	// 	utils.Tracef("eth_getlogs: receipt log: %+v", l)
	// }

	// return result, err
}

func (rpc *client) GetLogs(fromBlockNum, toBlockNum uint64, addresses []string, topics [][]string) ([]blockchain.IReceiptLog, error) {
	return nil, fmt.Errorf("not implemented")

	// filterParam := ethrpc.FilterParams{
	// 	FromBlock: "0x" + strconv.FormatUint(fromBlockNum, 16),
	// 	ToBlock:   "0x" + strconv.FormatUint(toBlockNum, 16),
	// 	Address:   addresses,
	// 	Topics:    topics,
	// }

	// logs, err := rpc.rpcImpl.EthGetLogs(filterParam)
	// if err != nil {
	// 	utils.Warnf("eth_getlogs: failed to retrieve logs: %v", err)
	// 	return nil, err
	// }

	// utils.Tracef("eth_getlogs: log count at block(%d - %d): %d", fromBlockNum, toBlockNum, len(logs))

	// var result []blockchain.IReceiptLog
	// for i := 0; i < len(logs); i++ {
	// 	l := logs[i]
	// 	result = append(result, blockchain.ReceiptLog{Log: &l})

	// 	utils.Tracef("eth_getlogs: receipt log: %+v", l)
	// }

	// return result, err
}

func (rpc *client) GetTransactionReceipt(txHash string) (blockchain.TransactionReceipt, error) {
	return nil, fmt.Errorf("not implemented")

	// receipt, err := rpc.rpcImpl.EthGetTransactionReceipt(txHash)
	// if err != nil {
	// 	return nil, err
	// }
	// if receipt == nil {
	// 	return nil, errors.New("nil receipt")
	// }

	// return &blockchain.EthereumTransactionReceipt{receipt}, err
}
