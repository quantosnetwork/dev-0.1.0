package tests

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"github.com/quantosnetwork/v0.1.0-dev/crypto/anon"
	"testing"
)

var SENTENCE_TO_ENCODE = "i am a quantos user yay me!"

var (
	testKey *[32]byte = new([32]byte)
)

func init() {
	rand.Read(testKey[:])
}

type TestAnon struct {
	authkey *[32]byte
	nonce   []byte
	encoded []byte
}

func TestAnonEncode(t *testing.T) {

	nonce := make([]byte, 8)
	f := func(pktNum uint64, in []byte) ([]byte, []byte, bool) {
		binary.BigEndian.PutUint64(nonce, pktNum)
		encoded, err := anon.EncodeNECS(testKey, nonce, in)
		if err != nil {
			return nil, nil, false
		}
		decoded, err := anon.NECSDecode(testKey, nonce, encoded)

		if err != nil {

			return nil, nil, false
		}

		return encoded, decoded, bytes.Compare(decoded, in) == 0
	}

	enc, _, istrue := f(0, []byte(SENTENCE_TO_ENCODE))

	if !istrue {
		t.Error("not true")
	} else {
		t.Logf("%x", enc)
	}

}
