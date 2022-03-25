package keygen

import (
	"github.com/cloudflare/circl/kem/kyber/kyber512"
	"lukechampine.com/frand"
)

/* KYBER KEYPAIR KEYGEN

 */

func (k *keys) KyberKeyGen() (*kyber512.PrivateKey, *kyber512.PublicKey) {
	pk, sk, err := kyber512.GenerateKeyPair(frand.Reader)
	if err != nil {
		panic(err)
	}
	k.Private = sk
	k.Public = pk
	return sk, pk
}
func (k *keys) GetKyberKeypair() *KyberKeyPair {
	kp := &KyberKeyPair{}
	kp.sk = k.Private.(*kyber512.PrivateKey)
	kp.pk = k.Public.(*kyber512.PublicKey)
	kp.RawPub, _ = kp.pk.MarshalBinary()
	kp.RawPriv, _ = kp.sk.MarshalBinary()
	return kp
}

type KyberKeyPair struct {
	sk      *kyber512.PrivateKey
	pk      *kyber512.PublicKey
	RawPub  []byte
	RawPriv []byte
}
