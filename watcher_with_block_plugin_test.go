package ethwatcher

import (
	"context"
	"testing"

	"github.com/rakshasa/ethwatcher/plugin"
	"github.com/rakshasa/ethwatcher/structs"
	"github.com/rakshasa/ethwatcher/utils"
	"github.com/sirupsen/logrus"
)

func TestNewBlockNumPlugin(t *testing.T) {
	utils.SetCategoryLogLevel(logrus.InfoLevel)

	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w := NewHttpBasedEthWatcher(api)

	utils.Printf("waiting for new block...")
	w.RegisterBlockPlugin(plugin.NewBlockNumPlugin(func(i uint64, b bool) {
		utils.Printf(">> found new block: %d, is removed: %t", i, b)
	}))

	w.RunTillExit(context.Background())
}

func TestSimpleBlockPlugin(t *testing.T) {
	api := "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
	w := NewHttpBasedEthWatcher(api)

	w.RegisterBlockPlugin(plugin.NewSimpleBlockPlugin(func(block *structs.RemovableBlock) {
		utils.Infof(">> %+v", block.Block)
	}))

	w.RunTillExit(context.Background())
}
