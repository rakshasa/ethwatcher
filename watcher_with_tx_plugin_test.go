package ethwatcher

import (
	"context"
	"testing"

	"github.com/rakshasa/ethwatcher/plugin"
	"github.com/rakshasa/ethwatcher/structs"
	"github.com/rakshasa/ethwatcher/utils"
)

func TestTxHashPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w := NewHttpBasedEthWatcher(api)

	w.RegisterTxPlugin(plugin.NewTxHashPlugin(func(txHash string, isRemoved bool) {
		utils.Printf("%s >> %s", txHash, isRemoved)
	}))

	w.RunTillExit(context.Background())
}

func TestTxPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w := NewHttpBasedEthWatcher(api)

	w.RegisterTxPlugin(plugin.NewTxPlugin(func(tx structs.RemovableTx) {
		utils.Printf(">> block: %d, txHash: %s", tx.GetBlockNumber(), tx.GetHash())
	}))

	w.RunTillExit(context.Background())
}
