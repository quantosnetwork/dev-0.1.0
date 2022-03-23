package common

import (
	"github.com/mr-tron/base58"
)

func EncodeBase58(data []byte) string {
	return base58.FastBase58Encoding(data)
}

func DecodeBase58(data string) []byte {
	b, err := base58.FastBase58Decoding(data)
	if err != nil {
		return nil
	}
	return b
}
