package rpc

import (
	"context"
	"errors"
	"strconv"

	// "github.com/ethereum/go-ethereum/ethclient"
	"github.com/onrik/ethrpc"
	"github.com/rakshasa/ethwatcher/blockchain"
	"github.com/rakshasa/ethwatcher/utils"
)

type Client struct {
	// client *ethclient.Client

	// Deprecate:
	rpcImpl *ethrpc.EthRPC
}

func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	// client, err := ethclient.DialContext(ctx, rawurl)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to dial ethereum api server: %v", err)
	// }

	return &Client{
		// client:  client,
		rpcImpl: ethrpc.New(rawurl),
	}, nil
}

func (rpc Client) NewFilter(addresses []string, topics []string) (string, error) {
	filterParams := ethrpc.FilterParams{
		FromBlock: "latest",
		ToBlock:   "latest",
		Address:   addresses,
		Topics:    [][]string{topics},
	}

	filterId, err := rpc.rpcImpl.EthNewFilter(filterParams)
	if err != nil {
		utils.Warnf("eth_newfilter: failed to create new filter: %v", err)
		return "", err
	}

	return filterId, err
}

func (rpc Client) GetBlockByNum(num uint64) (blockchain.Block, error) {
	b, err := rpc.rpcImpl.EthGetBlockByNumber(int(num), true)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("nil block")
	}

	return &blockchain.EthereumBlock{b}, err
}

func (rpc Client) GetBlockByNumWithoutTx(num uint64) (blockchain.Block, error) {
	b, err := rpc.rpcImpl.EthGetBlockByNumber(int(num), false)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("nil block")
	}

	return &blockchain.EthereumBlock{b}, err
}

func (rpc Client) GetCurrentBlockNum() (uint64, error) {
	num, err := rpc.rpcImpl.EthBlockNumber()
	return uint64(num), err
}

func (rpc Client) GetFilterChanges(filterId string) ([]blockchain.IReceiptLog, error) {
	logs, err := rpc.rpcImpl.EthGetFilterChanges(filterId)
	if err != nil {
		utils.Warnf("eth_getfilterchanges: failed to retrieve filter changes: %v", err)
		return nil, err
	}

	var result []blockchain.IReceiptLog
	for i := 0; i < len(logs); i++ {
		l := logs[i]
		result = append(result, blockchain.ReceiptLog{Log: &l})

		utils.Tracef("eth_getlogs: receipt log: %+v", l)
	}

	return result, err
}

func (rpc Client) GetLogs(
	fromBlockNum, toBlockNum uint64,
	addresses []string,
	topics []string,
) ([]blockchain.IReceiptLog, error) {

	filterParam := ethrpc.FilterParams{
		FromBlock: "0x" + strconv.FormatUint(fromBlockNum, 16),
		ToBlock:   "0x" + strconv.FormatUint(toBlockNum, 16),
		Address:   addresses,
		Topics:    [][]string{topics},
	}

	logs, err := rpc.rpcImpl.EthGetLogs(filterParam)
	if err != nil {
		utils.Warnf("eth_getlogs: failed to retrieve logs: %v", err)
		return nil, err
	}

	utils.Tracef("eth_getlogs: log count at block(%d - %d): %d", fromBlockNum, toBlockNum, len(logs))

	var result []blockchain.IReceiptLog
	for i := 0; i < len(logs); i++ {
		l := logs[i]
		result = append(result, blockchain.ReceiptLog{Log: &l})

		utils.Tracef("eth_getlogs: receipt log: %+v", l)
	}

	return result, err
}

func (rpc Client) GetTransactionReceipt(txHash string) (blockchain.TransactionReceipt, error) {
	receipt, err := rpc.rpcImpl.EthGetTransactionReceipt(txHash)
	if err != nil {
		return nil, err
	}
	if receipt == nil {
		return nil, errors.New("nil receipt")
	}

	return &blockchain.EthereumTransactionReceipt{receipt}, err
}
