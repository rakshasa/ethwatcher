package ethwatcher

import (
	"context"
	"fmt"
	"github.com/rakshasa/ethwatcher/plugin"
	"github.com/rakshasa/ethwatcher/structs"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestTxHashPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w := NewHttpBasedEthWatcher(context.Background(), api)

	w.RegisterTxPlugin(plugin.NewTxHashPlugin(func(txHash string, isRemoved bool) {
		fmt.Println(">>", txHash, isRemoved)
	}))

	w.RunTillExit()
}

func TestTxPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w := NewHttpBasedEthWatcher(context.Background(), api)

	w.RegisterTxPlugin(plugin.NewTxPlugin(func(tx structs.RemovableTx) {
		logrus.Printf(">> block: %d, txHash: %s", tx.GetBlockNumber(), tx.GetHash())
	}))

	w.RunTillExit()
}
