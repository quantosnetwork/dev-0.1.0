package hash

import (
	"encoding/hex"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

var edSuite = edwards25519.NewBlakeSHA256Ed25519()
var HashFn = edSuite.Hash()

type Hash struct {
	hash [32]byte
}

func NewHash(b []byte) *Hash {

	HashFn.Reset()
	HashFn.Write(b)
	s := HashFn.Sum(nil)
	h := &Hash{}
	copy(h.hash[:], s[:])
	return h
}

func (h *Hash) Bytes() []byte {
	return h.hash[:]
}

func (h *Hash) String() string {
	return hex.EncodeToString(h.hash[:])
}
