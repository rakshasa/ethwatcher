package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ethereumRPCEndpoint = "https://mainnet.infura.io/v3/"
	minBlockNumber      = uint64(14724100)
	maxBlockNumber      = minBlockNumber + 10000000
)

func TestClient(t *testing.T) {
	client, err := Dial(ethereumRPCEndpoint)
	if !t.NoError(err) {
		return
	}

	type testData struct {
		name string
		fn   func() error
	}

	tests := []testData{
		{
			"BlockNumber",
			func(assert *assert.Assertions) {
				v, err := client.BlockNumber(context.Background())
				if !assert.NoError(err) {
					return
				}
				if !assert.Greater(v, minBlockNumber) || !assert.Less(v, maxBlockNumber) {
					return
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(st *testing.T) {
			test.fn(assert.New(st))
		})
	}
}
