package ethwatcher

import (
	"context"
	"fmt"

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

func NewReceiptLogWatcher(
	api string,
	contractAddresses []string,
	interestedTopics []string,
	handler func(receiptLogs []blockchain.IReceiptLog) error,
	configs ...ReceiptLogWatcherConfig,
) *ReceiptLogWatcher {

	config := decideConfig(configs...)

	return &ReceiptLogWatcher{
		api:               api,
		contractAddresses: contractAddresses,
		interestedTopics:  interestedTopics,
		handler:           handler,
		config:            config,
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

type ReceiptLogWatcherConfig struct {
	RPCMaxRetry int
}

var defaultConfig = ReceiptLogWatcherConfig{
	RPCMaxRetry: 5,
}

func (w *ReceiptLogWatcher) Run(ctx context.Context) error {
	rpc := rpc.NewEthRPCWithRetry(w.api, w.config.RPCMaxRetry)

	filterId, err := rpc.NewFilter(w.contractAddresses, w.interestedTopics)
	if err != nil {
		return fmt.Errorf("failed to request new filter from api: %v", err)
	}

	for {
		utils.Debugf("polling eth filter changes...")

		// type queryResult struct {
		// 	logs []
		// }

		// TODO: Change to select wait for result.

		select {
		case <-ctx.Done():
			return nil
		default:
			logs, err := rpc.GetFilterChanges(filterId)
			if err != nil {
				return err
			}

			if len(logs) == 0 {
				// TODO: Properly handle this?
				continue
			}

			if err := w.handler(logs); err != nil {
				utils.Infof("err when handling receipt logs: %+v", logs)
				return fmt.Errorf("ethwatcher handler returns error: %s", err)
			}
		}
	}
}
