package anon

import (
	"crypto/rand"
	"github.com/quantosnetwork/v0.1.0-dev/common"
	"io"
	"log"
)

/*
	Encryption less confidentiality using the Chaffing and winnowing technology.

*/

const (
	NECS_ENLARGE_SIZE = HSize + RSize*256
)

var (
	Rand io.Reader = rand.Reader
)

func EncodeNECS(authKey *[32]byte, nonce, in []byte) ([]byte, error) {
	r := new([RSize]byte)
	var err error
	if _, err = Rand.Read(r[:]); err != nil {
		return nil, err
	}
	oaeped, err := Encode(r, in)
	if err != nil {
		return nil, err
	}
	out := append(Chaff(authKey, nonce, oaeped[:RSize]), oaeped[RSize:]...)
	common.SliceZero(oaeped[:RSize])
	return out, nil
}

func NECSDecode(authKey *[32]byte, nonce, in []byte) ([]byte, error) {

	winnowed, err := Winnow(
		authKey, nonce, in[:RSize*EnlargeFactor],
	)
	if err != nil {
		return nil, err
	}

	out, err := Decode(append(
		winnowed, in[RSize*EnlargeFactor:]...,
	))
	common.SliceZero(winnowed)
	if err != nil {
		return nil, err
	}
	log.Printf("%s", out)
	return out, nil
}
