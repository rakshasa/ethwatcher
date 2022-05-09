package rpc

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/rakshasa/ethwatcher/utils"
	"github.com/stretchr/testify/assert"
)

func newBigList(strs ...string) (list []big.Int) {
	list = make([]big.Int, len(strs))

	for idx, str := range strs {
		value, err := utils.NewBigFromHex(str, 8*32)
		if err != nil {
			utils.Fatalf("invalid hex string: %s", str)
			os.Exit(1)
		}

		list[idx] = *value
	}

	return
}

func newBig(str string) *big.Int {
	value, err := utils.NewBigFromHex(str, 8*32)
	if err != nil {
		utils.Fatalf("invalid hex string: %s", str)
		os.Exit(1)
	}

	return value
}

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
		"FilterLogs": func(ctx context.Context, assert *assert.Assertions) {
			v, err := client.FilterLogs(
				ctx,
				testBlockNumber,
				testBlockNumber,
				[]string{"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"},
				[][]string{[]string{"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"}},
			)
			if !assert.NoError(err) {
				return
			}
			if !assert.Len(v, 10) {
				return
			}

			assert.Equal(newBig("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), v[0].AddressAsBig())
			assert.Equal(
				newBigList(
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					"0x000000000000000000000000ee139899eadfbc5c25c7a1c6786d0367c82f7a1e",
					"0x0000000000000000000000000bf46ba06dc1d33c3bd80ff42497ebff13a88900",
				), v[0].TopicsAsBig())

			assert.Equal(newBig("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"), v[0].AddressAsBig())
			assert.Equal(
				newBigList(
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
					"0x0000000000000000000000008032eaede5c55f744387ca53aaf0499abcd783e5",
					"0x000000000000000000000000a5e9c917b4b821e4e0a5bbefce078ab6540d6b5e",
				), v[1].TopicsAsBig())
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
