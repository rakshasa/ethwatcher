package blockchain

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockChain interface {
	GetTokenBalance(tokenAddress, address string) decimal.Decimal
	GetTokenAllowance(tokenAddress, proxyAddress, address string) decimal.Decimal

	GetBlockNumber() (uint64, error)
	GetBlockByNumber(blockNumber uint64) (*Block, error)

	GetTransaction(ID string) (Transaction, error)
	GetTransactionReceipt(ID string) (TransactionReceipt, error)
	GetTransactionAndReceipt(ID string) (Transaction, TransactionReceipt, error)
}

type Transaction interface {
	GetBlockHash() string
	GetBlockNumber() uint64
	GetFrom() string
	GetGas() int
	GetGasPrice() big.Int
	GetHash() string
	GetTo() string
	GetValue() big.Int
}

type TransactionReceipt interface {
	GetResult() bool
	GetBlockNumber() uint64

	GetBlockHash() string
	GetTxHash() string
	GetTxIndex() int

	GetLogs() []IReceiptLog
}

type IReceiptLog interface {
	GetRemoved() bool
	GetLogIndex() int
	GetTransactionIndex() int
	GetTransactionHash() string
	GetBlockNum() int
	GetBlockHash() string
	GetAddress() string
	GetData() string
	GetTopics() []string
}

// compile time interface check
var _ BlockChain = &Ethereum{}

type EthereumTransaction struct {
	*ethtypes.Transaction
}

func (t *EthereumTransaction) GetFrom() string {
	return ""
	// return t.From
}

func (t *EthereumTransaction) GetGas() int {
	return int(t.Gas())
}

func (t *EthereumTransaction) GetGasPrice() big.Int {
	return *t.GasPrice()
}

func (t *EthereumTransaction) GetValue() big.Int {
	return *t.Value()
}

func (t *EthereumTransaction) GetTo() string {
	return t.To().String()
}

func (t *EthereumTransaction) GetHash() string {
	return t.Hash().String()
}
func (t *EthereumTransaction) GetBlockNumber() uint64 {
	return uint64(0)
	// return uint64(*t.BlockNumber)
}

type EthereumTransactionReceipt struct {
	// *ethtypes.TransactionReceipt
}

func (r *EthereumTransactionReceipt) GetLogs() (rst []IReceiptLog) {
	// for i := range r.Logs {
	// 	l := ReceiptLog{&r.Logs[i]}
	// 	rst = append(rst, l)
	// }

	return
}

func (r *EthereumTransactionReceipt) GetResult() bool {
	// res, err := strconv.ParseInt(r.Status, 0, 64)

	// if err != nil {
	// 	panic(err)
	// }

	// return res == 1

	return false
}

func (r *EthereumTransactionReceipt) GetBlockNumber() uint64 {
	return uint64(0)
	// return uint64(r.BlockNumber)
}

func (r *EthereumTransactionReceipt) GetBlockHash() string {
	return ""
	// return r.BlockHash
}

func (r *EthereumTransactionReceipt) GetTxHash() string {
	return ""
	// return r.TransactionHash
}

func (r *EthereumTransactionReceipt) GetTxIndex() int {
	return 0
	// return r.TransactionIndex
}

type ReceiptLog struct {
	*ethtypes.Log
}

func (log ReceiptLog) GetRemoved() bool {
	return log.Removed
}

func (log ReceiptLog) GetLogIndex() int {
	return int(log.Index)
}

func (log ReceiptLog) GetTransactionIndex() int {
	return int(log.TxIndex)
}

func (log ReceiptLog) GetTransactionHash() string {
	return log.TxHash.String()
}

func (log ReceiptLog) GetBlockNum() int {
	return int(log.BlockNumber)
}

func (log ReceiptLog) GetBlockHash() string {
	return log.BlockHash.String()
}

func (log ReceiptLog) GetAddress() string {
	return log.Address.String()
}

func (log ReceiptLog) GetData() string {
	return string(log.Data)
}

func (log ReceiptLog) GetTopics() (topics []string) {
	topics = make([]string, len(log.Topics))

	for idx, t := range log.Topics {
		topics[idx] = t.String()
	}

	return
}

type Ethereum struct {
	client       *ethclient.Client
	hybridExAddr string
}

func (e *Ethereum) EnableDebug(b bool) {
	// e.client.Debug = b
}

func (e *Ethereum) GetBlockByNumber(number uint64) (*Block, error) {
	return nil, fmt.Errorf("not implemented")

	// block, err := e.client.EthGetBlockByNumber(int(number), true)

	// if err != nil {
	// 	log.Errorf("get Block by Number failed %+v", err)
	// 	return nil, err
	// }

	// if block == nil {
	// 	log.Errorf("get Block by Number returns nil block for num: %d", number)
	// 	return nil, errors.New("get Block by Number returns nil block for num: " + strconv.Itoa(int(number)))
	// }

	// return &EthereumBlock{block}, nil
}

func (e *Ethereum) GetBlockNumber() (uint64, error) {
	return uint64(0), fmt.Errorf("not implemented")

	// number, err := e.client.EthBlockNumber()

	// if err != nil {
	// 	log.Errorf("GetBlockNumber failed, %v", err)
	// 	return 0, err
	// }

	// return uint64(number), nil
}

func (e *Ethereum) GetTransaction(ID string) (Transaction, error) {
	return nil, fmt.Errorf("not implemented")

	// tx, err := e.client.EthGetTransactionByHash(ID)

	// if err != nil {
	// 	log.Errorf("GetTransaction failed, %v", err)
	// 	return nil, err
	// }

	// return &EthereumTransaction{tx}, nil
}

func (e *Ethereum) GetTransactionReceipt(ID string) (TransactionReceipt, error) {
	return nil, fmt.Errorf("not implemented")

	// txReceipt, err := e.client.EthGetTransactionReceipt(ID)

	// if err != nil {
	// 	log.Errorf("GetTransactionReceipt failed, %v", err)
	// 	return nil, err
	// }

	// return &EthereumTransactionReceipt{txReceipt}, nil
}

func (e *Ethereum) GetTransactionAndReceipt(ID string) (Transaction, TransactionReceipt, error) {
	txReceiptChannel := make(chan TransactionReceipt)

	go func() {
		rec, _ := e.GetTransactionReceipt(ID)
		txReceiptChannel <- rec
	}()

	txInfoChannel := make(chan Transaction)
	go func() {
		tx, _ := e.GetTransaction(ID)
		txInfoChannel <- tx
	}()

	return <-txInfoChannel, <-txReceiptChannel, nil
}

func (e *Ethereum) GetTokenBalance(tokenAddress, address string) decimal.Decimal {
	return decimal.Decimal{}

	// res, err := e.client.EthCall(ethrpc.T{
	// 	To:   tokenAddress,
	// 	From: address,
	// 	Data: fmt.Sprintf("0x70a08231000000000000000000000000%s", without0xPrefix(address)),
	// }, "latest")

	// if err != nil {
	// 	panic(err)
	// }

	// return utils.StringToDecimal(res)
}

func without0xPrefix(address string) string {
	if address[:2] == "0x" {
		address = address[2:]
	}

	return address
}

func (e *Ethereum) GetTokenAllowance(tokenAddress, proxyAddress, address string) decimal.Decimal {
	return decimal.Decimal{}

	// res, err := e.client.EthCall(ethrpc.T{
	// 	To:   tokenAddress,
	// 	From: address,
	// 	Data: fmt.Sprintf("0xdd62ed3e000000000000000000000000%s000000000000000000000000%s", without0xPrefix(address), without0xPrefix(proxyAddress)),
	// }, "latest")

	// if err != nil {
	// 	panic(err)
	// }

	// return utils.StringToDecimal(res)
}

func (e *Ethereum) GetTransactionCount(address string) (int, error) {
	return 0, fmt.Errorf("not implemented")

	// return e.client.EthGetTransactionCount(address, "latest")
}
