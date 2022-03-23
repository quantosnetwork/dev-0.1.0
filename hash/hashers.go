package hash

import (
	"encoding/hex"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"lukechampine.com/blake3"
)

var edSuite = edwards25519.NewBlakeSHA256Ed25519()
var HashFn = edSuite.Hash()

type Hash struct {
	Hash []byte
}

func NewHash(b []byte) *Hash {

	HashFn.Reset()
	HashFn.Write(b)
	s := HashFn.Sum(nil)
	h := &Hash{}
	copy(h.Hash[:], s[:])
	return h
}

func (h *Hash) Bytes() []byte {
	return h.Hash[:]
}

func (h *Hash) String() string {
	return hex.EncodeToString(h.Hash[:])
}

func NewBlake3Hash(data []byte) []byte {
	h := blake3.New(256, nil) //unkeyed use func NewKeyedBlake3Hash() for Keyed
	_, err := h.Write(data)
	if err != nil {
		return nil
	}
	return h.Sum(nil)
}

func NewMultiDataBlake3Hash(data ...[]byte) *blake3.Hasher {
	h := blake3.New(256, nil)
	h.Write(data[0])
	for i := 1; i < len(data) && len(data) > 1; i++ {
		h.Write(data[i])
	}
	return h
}

func NewKeyedBlake3Hash(data []byte, key []byte) []byte {
	return nil
}

func GetBlake3Hasher(size int) *blake3.Hasher {
	return blake3.New(size, nil)
}
