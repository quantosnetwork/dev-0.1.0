package account

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mr-tron/base58"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/xof/blake2xb"
	"golang.org/x/crypto/scrypt"
	"hash/crc32"
	"io/ioutil"
	"log"
	"lukechampine.com/frand"
	"os"
	"path"
)

type IAddress interface {
	GenerateNewAddress() *address
	GetAddressFromFile(pubStr string) *Address
}

type AddressRepository struct {
	List map[string]bool
}

func (ar AddressRepository) New() *AddressRepository {
	r := new(AddressRepository)
	r.List = map[string]bool{}
	return r
}

func (ar *AddressRepository) Add(addr Address) {
	ar.List[addr.String()] = true
}

func (ar *AddressRepository) Disable(addr Address) {
	ar.List[addr.String()] = false
}

func (ar *AddressRepository) Has(addr Address) bool {
	return ar.List[addr.String()] == true
}

const AddressPrefix = "QBX"

type Address []byte

func (a Address) String() string {
	addr := fmt.Sprintf("QBX-%s", base58.Encode(a))
	return addr
}

func (a Address) Bytes() []byte {
	return a
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
	key, salt, err := DeriveKey(key, nil)
	if err != nil {
		return nil, err
	}
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
	ciphertext = append(ciphertext, salt...)
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
	dir, _ := os.Getwd()
	pat := path.Join(dir, "./../")
	err = ioutil.WriteFile(pat+"/data/.keys/"+hex.EncodeToString(a.PubBytes)+".wal", afw, 0600)
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

func GetAddressFromStorage(pub string, key string) (*DecryptedAddressFormat, error) {
	dir, _ := os.Getwd()
	pat := path.Join(dir, "./../")
	b, err := ioutil.ReadFile(pat + "/data/.keys/" + pub + ".wal")
	if err != nil {
		return nil, errors.New("invalid address, not found")
	}
	var fileFormat AddressFileFormat
	err = json.Unmarshal(b, &fileFormat)
	if err != nil {
		return nil, err
	}
	dec, err := Decrypt([]byte(key), fileFormat.EncryptedFile)
	if err != nil {
		return nil, errors.New("invalid key, address is locked")
	}
	var decAddr DecryptedAddressFormat
	_ = json.Unmarshal(dec, &decAddr)
	return &decAddr, nil

}

func CreateNewAddress(network byte, key string) (string, error) {
	a := new(address)
	addr := a.GenerateNewAddress(network)

	addrString := hex.EncodeToString(addr.PubBytes)
	err := addr.WriteAndLock([]byte(key))
	if err != nil {
		return "", err
	}
	return addrString, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := DeriveKey(key, salt)
	if err != nil {
		return nil, err
	}
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := frand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
