package util

import "encoding/base64"

func Byte2Base64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Base642Byte(h string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(h)
}
