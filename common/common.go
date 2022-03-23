package common

import (
	"golang.org/x/crypto/sha3"
)

// Zero each byte.
func SliceZero(data []byte) {
	for i := 0; i < len(data); i++ {
		data[i] = 0
	}
}

func Sha3(raw []byte) []byte {
	data := sha3.Sum256(raw)
	return data[:]
}
