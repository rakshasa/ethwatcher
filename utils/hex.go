package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

func Bytes2Hex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func Hex2Bytes(str string) []byte {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}

	if len(str)%2 == 1 {
		str = "0" + str
	}

	h, _ := hex.DecodeString(str)
	return h
}

// with prefix '0x'
func Bytes2HexP(bytes []byte) string {
	return "0x" + hex.EncodeToString(bytes)
}

func NewHashFromHexString(str string) (ethcommon.Hash, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}

	if len(str) != ethcommon.HashLength {
		return ethcommon.Hash{}, fmt.Errorf("invalid length")
	}

	b := bytes.ToLower([]byte(str))

	for _, c := range b {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return ethcommon.Hash{}, fmt.Errorf("not alphanumeric")
		}
	}

	return ethcommon.BytesToHash(b), nil
}
