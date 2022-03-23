package account

import (
	"bytes"
	"encoding/hex"
	"errors"
	dh "github.com/cloudflare/circl/dh/x25519"
	"github.com/cloudflare/circl/group"
	"github.com/cloudflare/circl/sign/ed25519"
	"io"
	"io/ioutil"
	"lukechampine.com/frand"
	"os"
)

type Keys struct {
	r         group.Group
	shared    map[dh.Key]dh.Key
	dhPub     dh.Key
	dhPriv    dh.Key
	PublicKey []byte
	privKey   []byte
	Sig       []byte
}

func NewDHKeys() (dh.Key, dh.Key) {
	var pub, secret dh.Key
	_, _ = io.ReadFull(frand.Reader, secret[:])
	dh.KeyGen(&pub, &secret)
	return pub, secret
}

func (k Keys) Shared(with dh.Key) {
	var shared dh.Key
	dh.Shared(&shared, &k.dhPriv, &with)
	k.shared[with] = shared
}

func hexStr2Key(k *dh.Key, s string) {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("Can't convert string to key")
	}
	copy(k[:], b)
}

func GenerateEd25519(id string, seed []byte) ed25519.PrivateKey {
	sk := ed25519.NewKeyFromSeed(seed)
	skb, _ := sk.MarshalBinary()

	WriteKeyForAccountID(id, seed, skb)
	return sk
}

func randSeed() []byte {
	b := frand.Entropy256()
	return b[:]
}

func WriteKeyForAccountID(accountID string, data []byte, priv []byte) {
	dir := "./data/.private/.s-" + accountID
	_ = ioutil.WriteFile(dir, data, 0600)
	_ = ioutil.WriteFile("./data/.keys/"+accountID+"-priv.key", priv, 0600)
}

func loadPrivateKey(accountID string, pubKey ed25519.PublicKey) (ed25519.PrivateKey, error) {
	sk, err := ioutil.ReadFile("./data/.keys/" + accountID + "-priv.key")
	s, err2 := ioutil.ReadFile("./data/.private/.s-" + accountID)
	if err != nil {
		if os.IsNotExist(err) {
			s := randSeed()
			sk = GenerateEd25519(accountID, s)

		}
	}
	if err2 != nil {
		if os.IsNotExist(err) {
			s := randSeed()
			sk = GenerateEd25519(accountID, s)

		}
	}

	sk2 := ed25519.NewKeyFromSeed(s)
	sk2b, _ := sk2.MarshalBinary()

	if bytes.Equal(sk, sk2b) {

		pub2 := sk2.Public()
		if pubKey.Equal(pub2) {
			return sk2, nil
		} else {
			return nil, errors.New("private key is not valid")
		}
	}
	return nil, errors.New("private key is not valid")

}

func (k Keys) signPublic(priv ed25519.PrivateKey) []byte {
	// self-sign pubkey
	sig := ed25519.Sign(priv, k.PublicKey)
	return sig
}

type LoadedKeys struct {
	Priv       ed25519.PrivateKey
	Pub        ed25519.PublicKey
	PubKeySign []byte
}

func (k Keys) GetLoadedKeys(accID string) *LoadedKeys {
	sk, err := loadPrivateKey(accID, k.PublicKey)
	if err != nil {
		panic(err)
	}
	return &LoadedKeys{
		Priv:       sk,
		Pub:        sk.Public().(ed25519.PublicKey),
		PubKeySign: k.signPublic(sk),
	}
}

func initializeKeyPair(accID string) (ed25519.PrivateKey, ed25519.PublicKey) {
	var K Keys
	l := K.GetLoadedKeys(accID)
	return l.Priv, l.Pub
}
