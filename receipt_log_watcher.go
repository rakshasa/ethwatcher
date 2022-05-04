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
	config            receiptLogWatcherConfig
}

type receiptLogWatcherConfig struct {
	blockStepSize   uint64
	pollingInterval time.Duration
	rpcMaxRetries   int
	useFilter       bool
}

func WithRPCMaxRetries(retries int) func(*receiptLogWatcherConfig) {
	return func(config *receiptLogWatcherConfig) {
		config.rpcMaxRetries = retries
	}
}

func NewReceiptLogWatcher(
	api string,
	contractAddresses []string,
	interestedTopics []string,
	handler func(receiptLogs []blockchain.IReceiptLog) error,
	options ...func(*receiptLogWatcherConfig),
) *ReceiptLogWatcher {
	config := receiptLogWatcherConfig{
		blockStepSize:   50,
		pollingInterval: 5 * time.Second,
		rpcMaxRetries:   5,
		useFilter:       false,
	}

	for _, optFn := range options {
		optFn(&config)
	}

	return &ReceiptLogWatcher{
		api:               api,
		contractAddresses: contractAddresses,
		interestedTopics:  interestedTopics,
		handler:           handler,
		config:            config,
	}
}

func (w *ReceiptLogWatcher) Run(ctx context.Context) error {
	if w.config.rpcMaxRetries < 0 {
		return fmt.Errorf("invalid RPCMaxRetries value")
	}

	rpc := rpc.NewEthRPCWithRetry(w.api, uint(w.config.rpcMaxRetries))

	var err error
	var filterId string
	var prevBlockNum, nextBlockNum uint64

	if w.config.useFilter {
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

		if w.config.useFilter {
			logs, err = rpc.GetFilterChanges(filterId)
			if err != nil {
				return err
			}
		} else {
			nextBlockNum, err = rpc.GetCurrentBlockNum()
			if err != nil {
				return err
			}

			if nextBlockNum-prevBlockNum > w.config.blockStepSize {
				nextBlockNum = prevBlockNum + w.config.blockStepSize
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

		if w.config.useFilter {
		} else {
			prevBlockNum = nextBlockNum
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(w.config.pollingInterval - time.Since(prevTime)):
		}
	}
}
