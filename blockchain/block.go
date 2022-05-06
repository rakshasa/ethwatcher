package blockchain

import (
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	block ethtypes.Block
}

func NewBlock(block *ethtypes.Block) *Block {
	return &Block{*block}
}

func (block *Block) Hash() string {
	return block.block.Hash().String()
}

func (block *Block) ParentHash() string {
	return block.block.ParentHash().String()
}

func (block *Block) Number() uint64 {
	return block.block.NumberU64()
}

func (block *Block) Time() time.Time {
	return time.Unix(int64(block.block.Time()), int64(0))
}

func (block *Block) TimeUint64() uint64 {
	return block.block.Time()
}

func (block *Block) GetTransactions() []Transaction {
	return nil

	// txs := make([]Transaction, 0, 20)

	// for i := range block.Block.Transactions {
	// 	tx := block.Block.Transactions[i]
	// 	txs = append(txs, &EthereumTransaction{&tx})
	// }

	// return txs
}
