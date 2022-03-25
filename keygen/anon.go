package keygen

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/sign/anon"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
)

/*
Kyber based anonymous encryption
*/

var X []kyber.Point
var y []kyber.Scalar

type EncryptedData struct {
	suite      *edwards25519.SuiteEd25519
	X          []kyber.Point
	y          []kyber.Scalar
	cipherText []byte
}

type AnonEncrypt interface {
	GenEncryptionKeys(keyNum int) (s *edwards25519.SuiteEd25519, X []kyber.Point, y []kyber.Scalar)
	Encrypt(X []kyber.Point, suite *edwards25519.SuiteEd25519, m []byte) []byte
}

func (k *keys) GenEncryptionKeys(keyNum int) (s *edwards25519.SuiteEd25519, X []kyber.Point, y []kyber.Scalar) {
	s = edwards25519.NewBlakeSHA256Ed25519WithRand(blake2xb.New(nil))
	X = make([]kyber.Point, keyNum)
	y = make([]kyber.Scalar, keyNum)
	for i := 0; i < keyNum; i++ {
		y[i] = s.Scalar().Pick(s.RandomStream())
		X[i] = s.Point().Mul(y[i], nil)
	}
	return
}

func (k *keys) Encrypt(X []kyber.Point, suite *edwards25519.SuiteEd25519, m []byte) []byte {
	C := anon.Encrypt(suite, m, anon.Set(X))
	return C
}

func Encrypt(k *keys, keyNum int, msg []byte) *EncryptedData {
	e := &EncryptedData{}
	e.suite, e.X, e.y = k.GenEncryptionKeys(keyNum)
	C := k.Encrypt(e.X, e.suite, msg)
	e.cipherText = C
	return e
}

func (e *EncryptedData) GetSuite() *edwards25519.SuiteEd25519 {
	return e.suite
}

func (e *EncryptedData) GetCipherText() []byte {
	return e.cipherText
}

func (e *EncryptedData) GetX() []kyber.Point {
	return e.X
}

func Decrypt(e *EncryptedData, sks ...any) [][]byte {
	buf := make([][]byte, len(sks))
	for i := 0; i < len(sks); i++ {
		buf[i], _ = anon.Decrypt(e.GetSuite(), e.GetCipherText(), anon.Set(e.GetX()), i, sks[i].(kyber.Scalar))
	}
	return buf
}
