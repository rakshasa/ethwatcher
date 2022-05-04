package ethwatcher

import (
	"context"
	"fmt"
	"time"

	"github.com/rakshasa/ethwatcher/blockchain"
	"github.com/rakshasa/ethwatcher/rpc"
	"github.com/rakshasa/ethwatcher/utils"
)

type ReceiptLogWatcher struct {
	api               string
	contractAddresses []string
	interestedTopics  []string
	handler           func(receiptLogs []blockchain.IReceiptLog) error
	config            ReceiptLogWatcherConfig
}

type ReceiptLogWatcherConfig struct {
	BlockStepSize   uint64
	PollingInterval time.Duration
	RPCMaxRetry     int
	UseFilter       bool
}

var defaultConfig = ReceiptLogWatcherConfig{
	BlockStepSize:   50,
	PollingInterval: 5 * time.Second,
	RPCMaxRetry:     5,
	UseFilter:       false,
}

func NewReceiptLogWatcher(
	api string,
	contractAddresses []string,
	interestedTopics []string,
	handler func(receiptLogs []blockchain.IReceiptLog) error,
	configs ...ReceiptLogWatcherConfig,
) *ReceiptLogWatcher {
	return &ReceiptLogWatcher{
		api:               api,
		contractAddresses: contractAddresses,
		interestedTopics:  interestedTopics,
		handler:           handler,
		config:            decideConfig(configs...),
	}
}

func decideConfig(configs ...ReceiptLogWatcherConfig) ReceiptLogWatcherConfig {
	if len(configs) == 0 {
		return defaultConfig
	}

	config := configs[0]

	if config.RPCMaxRetry <= 0 {
		config.RPCMaxRetry = defaultConfig.RPCMaxRetry
	}

	return config
}

func (w *ReceiptLogWatcher) Run(ctx context.Context) error {
	rpc := rpc.NewEthRPCWithRetry(w.api, w.config.RPCMaxRetry)

	var err error
	var filterId string
	var prevBlockNum, nextBlockNum uint64

	if w.config.UseFilter {
		if filterId, err = rpc.NewFilter(w.contractAddresses, w.interestedTopics); err != nil {
			return fmt.Errorf("failed to request new filter from api: %v", err)
		}
	} else {
		if prevBlockNum, err = rpc.GetCurrentBlockNum(); err != nil {
			return err
		}
	}

	for {
		utils.Debugf("polling eth receipt log changes...")

		prevTime := time.Now()

		var logs []blockchain.IReceiptLog

		if w.config.UseFilter {
			logs, err = rpc.GetFilterChanges(filterId)
			if err != nil {
				return err
			}
		} else {
			nextBlockNum, err = rpc.GetCurrentBlockNum()
			if err != nil {
				return err
			}

			if nextBlockNum-prevBlockNum > w.config.BlockStepSize {
				nextBlockNum = prevBlockNum + w.config.BlockStepSize
			}

			logs, err = rpc.GetLogs(prevBlockNum, nextBlockNum, w.contractAddresses, w.interestedTopics)
			if err != nil {
				return err
			}
		}

		// TODO: Properly handle empty results.

		if len(logs) != 0 {
			if err := w.handler(logs); err != nil {
				utils.Infof("error handling receipt logs: %+v", logs)
				return fmt.Errorf("could not handle receipt logs: %v", err)
			}
		}

		if w.config.UseFilter {
		} else {
			prevBlockNum = nextBlockNum
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(w.config.PollingInterval - time.Since(prevTime)):
		}
	}
}
