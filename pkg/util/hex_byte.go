package util

import (
	"encoding/hex"
)

func Byte2Hex(b []byte) string {
	return hex.EncodeToString(b)
}

func Hex2Byte(h string) ([]byte, error) {
	return hex.DecodeString(h)
}
