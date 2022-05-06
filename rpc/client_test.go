package rpc

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	skipRpcTests := os.Getenv("SKIP_RPC_TESTS")
	if skipRpcTests == "yes" {
		return
	}

	rpcEndpoint := os.Getenv("ETHEREUM_RPC_ENDPOINT")
	if !assert.NotEmpty(t, rpcEndpoint, "missing ETHEREUM_RPC_ENDPOINT environment variable") {
		return
	}

	client, err := Dial(rpcEndpoint)
	if !assert.NoError(t, err) {
		return
	}

	minBlockNumber := uint64(14724100)
	maxBlockNumber := minBlockNumber + 10000000
	testBlockNumber := uint64(14724302)

	tests := map[string]func(context.Context, *assert.Assertions){
		"BlockNumber": func(ctx context.Context, assert *assert.Assertions) {
			v, err := client.BlockNumber(ctx)
			if !assert.NoError(err) {
				return
			}

			assert.Less(minBlockNumber, v)
			assert.Greater(maxBlockNumber, v)
		},
		"BlockByNumber": func(ctx context.Context, assert *assert.Assertions) {
			v, err := client.BlockByNumber(ctx, testBlockNumber)
			if !assert.NoError(err) {
				return
			}

			assert.Equal("0xed324400ca83f5ecd257fd6e2626e87fe5cb641040b978fd6c7dd90b051d817c", v.Hash())
			assert.Equal("0x6da27232a35dce5ecc5659f95930b295054873794dc6f17a153d059880c7e8da", v.ParentHash())
			assert.Equal(testBlockNumber, v.Number())
			assert.Equal(uint64(0x62753d34), v.TimeUint64())
		},
	}

	for name, fn := range tests {
		t.Run(name, func(st *testing.T) {
			ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFn()

			fn(ctx, assert.New(st))
		})
	}
}
