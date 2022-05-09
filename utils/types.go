package utils

import (
	"fmt"
	"math/big"
	"strings"
)

func NewBigFromHex(str string, maxBits int) (*big.Int, error) {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}

	value, ok := new(big.Int).SetString(str, 16)
	if !ok {
		return nil, fmt.Errorf("invalid characters")
	}
	if value.BitLen() > maxBits {
		return nil, fmt.Errorf("invalid length")
	}

	return value, nil
}
