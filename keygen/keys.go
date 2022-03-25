package keygen

import (
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

type Type uint32

const (
	NULL_TYPE Type = iota
	KYBER
	DILITHIUM
	Ed25519
	ANON
)

var TypeToString map[Type]string = map[Type]string{
	NULL_TYPE: "NULL",
	KYBER:     "KYBER",
	DILITHIUM: "DILITHIUM",
	Ed25519:   "ED25519",
	ANON:      "ANON",
}

type KeyGen interface {
	NewKeyPair(keyType Type)
	GetPrivate() *privateKey
	GetPublic() *publicKey
}

type kyberPrivateKey kyber.Scalar
type kyberPublicKey kyber.Point

type privateKey struct{}
type publicKey struct{}

type keys struct {
	KeyType Type
	Private any
	Public  any
	Suite   *edwards25519.SuiteEd25519
}

func (k *keys) NewKeyPair(keyType Type) {
	switch keyType {
	case NULL_TYPE:
		break
	case KYBER:
		k.KeyType = KYBER
		kp := k.GetKyberKeypair()
		k.Private = kp.RawPriv
		k.Public = kp.RawPub
		break
	case ANON:
		// ANON purpose is for multi encryption only
		k.KeyType = ANON
		k.Suite, k.Private, k.Public = k.GenEncryptionKeys(3)
		break
	case DILITHIUM:
		break
	case Ed25519:
		break
	}
}

func GetKeyPair(keyType Type) *keys {
	k := &keys{}
	k.NewKeyPair(keyType)
	return k
}
