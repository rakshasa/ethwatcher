package ethwatcher

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/rakshasa/ethwatcher/blockchain"
	"github.com/rakshasa/ethwatcher/rpc"
	"github.com/rakshasa/ethwatcher/utils"
)

const (
	MaxBlockStepSize = 3500
)

type ReceiptLogWatcher struct {
	api               string
	contractAddresses []string
	interestedTopics  [][]string
	handler           func(receiptLogs []blockchain.Log) error
	config            receiptLogWatcherConfig
}

type receiptLogWatcherConfig struct {
	blockStepSize   uint
	pollingInterval time.Duration
	rpcMaxRetries   int
	useFilter       bool
}

func WithBlockStepSize(stepSize uint) func(*receiptLogWatcherConfig) {
	return func(config *receiptLogWatcherConfig) {
		config.blockStepSize = stepSize
	}
}

func WithPollingInterval(interval time.Duration) func(*receiptLogWatcherConfig) {
	return func(config *receiptLogWatcherConfig) {
		config.pollingInterval = interval
	}
}

func WithRPCMaxRetries(retries int) func(*receiptLogWatcherConfig) {
	return func(config *receiptLogWatcherConfig) {
		config.rpcMaxRetries = retries
	}
}

func WithUseFilter(use bool) func(*receiptLogWatcherConfig) {
	return func(config *receiptLogWatcherConfig) {
		config.useFilter = use
	}
}

func NewReceiptLogWatcher(api string, handler func(receiptLogs []blockchain.Log) error, options ...func(*receiptLogWatcherConfig)) *ReceiptLogWatcher {
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
		contractAddresses: []string{},
		interestedTopics:  [][]string{},
		handler:           handler,
		config:            config,
	}
}

func (w *ReceiptLogWatcher) SetFilter(contractAddresses []string, interestedTopics [][]string) {
	w.contractAddresses = contractAddresses
	w.interestedTopics = interestedTopics
}

func (w *ReceiptLogWatcher) Run(ctx context.Context) error {
	rpc, err := rpc.DialContextWithRetry(ctx, w.api, uint(w.config.rpcMaxRetries))
	if err != nil {
		return fmt.Errorf("failed to dial ethereum rpc node: %v", err)
	}

	prevBlockNum, err := rpc.BlockNumber(ctx)
	if err != nil {
		return err
	}

	return w.run(ctx, rpc, prevBlockNum, math.MaxUint64)
}

func (w *ReceiptLogWatcher) RunBetweenBlockNumbers(ctx context.Context, firstBlockNum, lastBlockNum uint64) error {
	rpc, err := rpc.DialContextWithRetry(ctx, w.api, uint(w.config.rpcMaxRetries))
	if err != nil {
		return fmt.Errorf("failed to dial ethereum rpc node: %v", err)
	}

	if firstBlockNum > 1 {
		firstBlockNum--
	}

	w.config.blockStepSize = MaxBlockStepSize

	return w.run(ctx, rpc, firstBlockNum, lastBlockNum)
}

func (w *ReceiptLogWatcher) run(ctx context.Context, rpc rpc.Client, prevBlockNum, lastBlockNum uint64) error {
	// TODO: These checks should be done during initialization.
	if w.config.blockStepSize > MaxBlockStepSize {
		return fmt.Errorf("invalid BlockStepSize value")
	}
	if w.config.pollingInterval < 0 {
		return fmt.Errorf("invalid PollingInterval value")
	}
	if w.config.rpcMaxRetries < 0 {
		return fmt.Errorf("invalid RPCMaxRetries value")
	}

	var err error
	var filterId string
	var nextBlockNum uint64

	if w.config.useFilter {
		if filterId, err = rpc.NewFilter(w.contractAddresses, w.interestedTopics); err != nil {
			return fmt.Errorf("failed to request new filter from api: %v", err)
		}
	}

	for {
		utils.Debugf("polling eth receipt log changes...")

		prevTime := time.Now()

		var logs []blockchain.Log

		if w.config.useFilter {
			logs, err = rpc.GetFilterChanges(filterId)
			if err != nil {
				return err
			}
		} else {
			nextBlockNum, err = rpc.BlockNumber(ctx)
			if err != nil {
				return err
			}

			if nextBlockNum > lastBlockNum {
				nextBlockNum = lastBlockNum
			}

			if nextBlockNum > prevBlockNum {
				if nextBlockNum-prevBlockNum > uint64(w.config.blockStepSize) {
					nextBlockNum = prevBlockNum + uint64(w.config.blockStepSize)
				}

				logs, err = rpc.FilterLogs(ctx, prevBlockNum+1, nextBlockNum, w.contractAddresses, w.interestedTopics)
				if err != nil {
					return err
				}
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

			if nextBlockNum < lastBlockNum {
				continue
			}
			if nextBlockNum >= lastBlockNum {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(w.config.pollingInterval - time.Since(prevTime)):
		}
	}
}
