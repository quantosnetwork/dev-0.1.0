package anon

import (
	"crypto/subtle"
	"errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/salsa20"
	"log"
)

// ANON takes 128-bit random r, data M to be encoded and produce the
// package PKG:
//
//     PKG = P1 || P2
//      P1 = Salsa20(key=r, nonce=0x00, 0x00) XOR (M || BLAKE2b(r || M))
//      P2 = BLAKE2b(P1) XOR r || BLAKE3(P1) XOR r

const (
	HSize = 32

	RSize = 16
)

var (
	fakeNonce []byte = make([]byte, 8)
)

func Encode(r *[RSize]byte, in []byte) ([]byte, error) {
	out := make([]byte, len(in)+HSize+RSize)
	copy(out, in)
	hb, _ := blake2b.New256(nil)
	hb.Write(r[:])
	hb.Write(in)
	copy(out[len(in):], hb.Sum(nil))
	salsaKey := new([32]byte)
	copy(salsaKey[:], r[:])
	salsa20.XORKeyStream(out, out, fakeNonce, salsaKey)
	hb.Reset()
	hb.Write(out[:len(in)+32])
	for i, b := range hb.Sum(nil)[:RSize] {
		out[len(in)+32+i] = b ^ r[i]
	}
	return out, nil
}

func Decode(in []byte) ([]byte, error) {
	log.Println(len(in))
	if len(in) < HSize+RSize {
		return nil, errors.New("too small input buffer")
	}
	hasher, _ := blake2b.New256(nil)

	hasher.Write(in[:len(in)-RSize])
	salsaKey := new([32]byte)
	for i, b := range hasher.Sum(nil)[:RSize] {
		salsaKey[i] = b ^ in[len(in)-RSize+i]
	}
	hasher.Reset()
	hasher.Write(salsaKey[:RSize])
	out := make([]byte, len(in)-RSize)
	log.Printf("salsakey: %x, fakeNonce: %x, lenInput: %d, output: %x", salsaKey, fakeNonce, len(in),
		out)
	salsa20.XORKeyStream(out, in[:len(in)-RSize], fakeNonce, salsaKey)
	hasher.Write(out[:len(out)-HSize])
	if subtle.ConstantTimeCompare(hasher.Sum(nil), out[len(out)-HSize:]) != 1 {
		return nil, errors.New("invalid checksum")
	}
	return out[:len(out)-HSize], nil

}
