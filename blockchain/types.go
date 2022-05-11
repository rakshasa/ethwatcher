package blockchain

import (
	"encoding/binary"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type Topic struct {
	ethcommon.Hash
}

func (t Topic) AddressHex() (string, bool) {
	for _, c := range t.Hash[0 : len(t.Hash)-ethcommon.AddressLength] {
		if c != 0 {
			return "", false
		}
	}

	return fmt.Sprintf("0x%040x", t.Hash[len(t.Hash)-ethcommon.AddressLength:]), true
}

func (t Topic) Hex() string {
	return fmt.Sprintf("0x%064x", t.Hash[0:len(t.Hash)])
}

func (t Topic) String() string {
	return t.Hex()
}

func (t Topic) Uint8() (uint8, bool) {
	for _, c := range t.Hash[0 : len(t.Hash)-1] {
		if c != 0 {
			return uint8(0), false
		}
	}

	return uint8(t.Hash[len(t.Hash)-1]), true
}

func (t Topic) Uint32() (uint32, bool) {
	for _, c := range t.Hash[0 : len(t.Hash)-4] {
		if c != 0 {
			return uint32(0), false
		}
	}

	return binary.BigEndian.Uint32(t.Hash[len(t.Hash)-4:]), true
}

func (t Topic) Uint64() (uint64, bool) {
	for _, c := range t.Hash[0 : len(t.Hash)-8] {
		if c != 0 {
			return uint64(0), false
		}
	}

	return binary.BigEndian.Uint64(t.Hash[len(t.Hash)-8:]), true
}
