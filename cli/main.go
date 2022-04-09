package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/rakshasa/ethwatcher"
	"github.com/rakshasa/ethwatcher/blockchain"
	"github.com/rakshasa/ethwatcher/plugin"
	"github.com/rakshasa/ethwatcher/rpc"
	"github.com/rakshasa/ethwatcher/utils"
	"github.com/spf13/cobra"
)

const (
	api = "https://mainnet.infura.io/v3/19d753b2600445e292d54b1ef58d4df4"
)

var contractAdx string
var eventSigs []string
var blockBackoff int

func main() {
	rootCMD.AddCommand(blockNumCMD)
	rootCMD.AddCommand(usdtTransferCMD)

	contractEventListenerCMD.Flags().StringVarP(&contractAdx, "contract", "c", "", "contract address listen to")
	contractEventListenerCMD.MarkFlagRequired("contract")
	contractEventListenerCMD.Flags().StringArrayVarP(&eventSigs, "events", "e", []string{}, "signatures of events we are interested in")
	contractEventListenerCMD.MarkFlagRequired("events")
	contractEventListenerCMD.Flags().IntVar(&blockBackoff, "block-backoff", 0, "how many blocks we go back")
	rootCMD.AddCommand(contractEventListenerCMD)

	if err := rootCMD.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCMD = &cobra.Command{
	Use:   "ethwatcher",
	Short: "ethwatcher makes getting updates from Ethereum easier",
}

var blockNumCMD = &cobra.Command{
	Use:   "new-block-number",
	Short: "Print number of new block",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(cmd.Context())

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)

		w := ethwatcher.NewHttpBasedEthWatcher(api)

		utils.Printf("waiting for new block...")
		w.RegisterBlockPlugin(plugin.NewBlockNumPlugin(func(i uint64, b bool) {
			utils.Printf(">> found new block: %d, is removed: %t", i, b)
		}))

		go func() {
			<-c
			cancel()
		}()

		if err := w.RunTillExit(ctx); err != nil {
			utils.Printf("exit with err: %s", err)
			return
		}

		utils.Infof("exit")
	},
}

var usdtTransferCMD = &cobra.Command{
	Use:   "usdt-transfer",
	Short: "Show Transfer Event of USDT",
	Run: func(cmd *cobra.Command, args []string) {
		usdtContractAdx := "0xdac17f958d2ee523a2206206994597c13d831ec7"

		// Transfer
		topicsInterestedIn := []string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}

		handler := func(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {

			if from != to {
				utils.Infof("See new USDT Transfer at blockRange: %d -> %d, count: %2d", from, to, len(receiptLogs))
			} else {
				utils.Infof("See new USDT Transfer at block: %d, count: %2d", from, len(receiptLogs))
			}

			for _, log := range receiptLogs {
				utils.Infof("  >> tx: https://etherscan.io/tx/%s", log.GetTransactionHash())
			}

			fmt.Println("  ")

			return nil
		}

		receiptLogWatcher := ethwatcher.NewReceiptLogWatcher(
			api,
			-1,
			[]string{usdtContractAdx},
			topicsInterestedIn,
			handler,
			ethwatcher.ReceiptLogWatcherConfig{
				StepSizeForBigLag:               5,
				IntervalForPollingNewBlockInSec: 5,
				RPCMaxRetry:                     3,
				ReturnForBlockWithNoReceiptLog:  true,
			},
		)

		receiptLogWatcher.Run(cmd.Context())
	},
}

var contractEventListenerCMD = &cobra.Command{
	Use:   "contract-event-listener",
	Short: "listen and print events from contract",
	Example: `
  listen to Transfer & Approve events from Multi-Collateral-DAI
  
  /bin/ethwatcher contract-event-listener \
    --block-backoff 100
    --contract 0x6b175474e89094c44da98b954eedeac495271d0f \
    --events 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef`,
	Run: func(cmd *cobra.Command, args []string) {

		handler := func(from, to int, receiptLogs []blockchain.IReceiptLog, isUpToHighestBlock bool) error {

			if from != to {
				utils.Infof("# of interested events at block(%d->%d): %d", from, to, len(receiptLogs))
			} else {
				utils.Infof("# of interested events at block(%d): %d", from, len(receiptLogs))
			}

			for _, log := range receiptLogs {
				utils.Infof("  >> tx: https://etherscan.io/tx/%s", log.GetTransactionHash())
			}

			fmt.Println("  ")

			return nil
		}

		startBlockNum := -1
		if blockBackoff > 0 {
			rpc := rpc.NewEthRPCWithRetry(api, 3)
			curBlockNum, err := rpc.GetCurrentBlockNum()
			if err == nil {
				startBlockNum = int(curBlockNum) - blockBackoff

				if startBlockNum > 0 {
					utils.Infof("--block-backoff activated, we start from block: %d (= %d - %d)",
						startBlockNum, curBlockNum, blockBackoff)
				}
			}
		}

		receiptLogWatcher := ethwatcher.NewReceiptLogWatcher(
			api,
			startBlockNum,
			[]string{contractAdx},
			eventSigs,
			handler,
			ethwatcher.ReceiptLogWatcherConfig{
				StepSizeForBigLag:               5,
				IntervalForPollingNewBlockInSec: 5,
				RPCMaxRetry:                     3,
				ReturnForBlockWithNoReceiptLog:  true,
			},
		)

		receiptLogWatcher.Run(cmd.Context())
	},
}
