package account

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
	"hash/crc32"
	"io/ioutil"
	"log"
	"lukechampine.com/frand"
)

type Address interface {
	GenerateNewAddress() *address
}

type address struct {
	Network  byte
	PubBytes []byte
	Private  []byte
	Suite    *edwards25519.SuiteEd25519
}

func (a *address) GenerateNewAddress(network byte) *address {
	seed := make([]byte, 32)
	frand.Read(seed)
	suite := edwards25519.NewBlakeSHA256Ed25519WithRand(blake2xb.New(seed))

	priv := suite.Scalar().Pick(suite.RandomStream())
	pub := suite.Point().Mul(priv, nil)
	log.Println(pub.String())
	pubBytes, _ := pub.MarshalBinary()
	a.PubBytes = pubBytes
	a.Private, _ = priv.MarshalBinary()
	a.Suite = suite
	a.Network = network
	return a

}

func (a *address) PrepareForWriting() []byte {
	d := &DecryptedAddressFormat{
		Network:   a.Network,
		PrivKey:   a.Private,
		PubKey:    a.PubBytes,
		Signature: nil,
		Suite:     a.Suite.String(),
	}
	dBytes, _ := json.Marshal(d)
	return dBytes
}

func (a *address) EncryptFormat(key []byte, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = frand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil

}

func (a *address) WriteAndLock(key []byte) error {
	d := a.PrepareForWriting()
	e, err := a.EncryptFormat(key, d)
	if err != nil {
		panic(err)
	}
	af := &AddressFileFormat{
		EncryptedFile: e,
		Chksum:        crc32.ChecksumIEEE(e),
	}
	afw, _ := json.Marshal(af)
	err = ioutil.WriteFile("./data/.keys/"+hex.EncodeToString(a.PubBytes)+".wal", afw, 0600)
	if err != nil {
		return err
	}
	return nil
}

type AddressFileFormat struct {
	EncryptedFile []byte
	Chksum        uint32
}

type DecryptedAddressFormat struct {
	Network   byte
	PrivKey   []byte
	PubKey    []byte
	Signature []byte
	Suite     string
}

func CreateNewAddress(network byte, key string) string {
	a := new(address)
	addr := a.GenerateNewAddress(0x1e)
	spew.Dump(addr)
	addrString := hex.EncodeToString(addr.PubBytes)
	err := addr.WriteAndLock([]byte(key))
	if err != nil {
		panic(err)
	}
	return addrString
}
