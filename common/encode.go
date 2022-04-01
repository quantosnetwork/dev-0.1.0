package common

import (
	"github.com/mr-tron/base58"
)

type Base58Encoder interface {
	DoEncodeBase58(data []byte) string
	DoDecodeBase58(data string) []byte
}

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

type Encoder struct {
	Base58Encoder
}

func (b *Encoder) DoEncodeBase58(data []byte) string {
	return EncodeBase58(data)
}

func (b *Encoder) DoDecodeBase58(data string) []byte {
	return DecodeBase58(data)
}
