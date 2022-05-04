#!/bin/bash

set -eux

project_root="$(cd "$(git -C "$( dirname "${BASH_SOURCE[0]}" )" rev-parse --show-toplevel)" && pwd)"; readonly project_root


dependencies=(
  "github.com/btcsuite/btcd/btcec@v0.0.0-20190115013929-ed77733ec07d"
  "github.com/ethereum/go-ethereum/ethclient@v1.10.17"
  "golang.org/x/crypto@latest"
)

cd "${project_root}"

rm -rf ./go.mod ./go.sum

go clean -cache
go mod init github.com/rakshasa/ethwatcher

for dep in "${dependencies[@]}"; do
  go get "${dep}"
done

go mod tidy
