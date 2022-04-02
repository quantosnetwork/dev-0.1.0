package super

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

func marshalRSAPrivate(priv *rsa.PrivateKey) string {
	return string(pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}))
}

func generateKey() (string, string, error) {
	reader := rand.Reader
	bitSize := 2048

	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		return "", "", err
	}

	pub, err := ssh.NewPublicKey(key.Public())
	if err != nil {
		return "", "", err
	}
	pubKeyStr := string(ssh.MarshalAuthorizedKey(pub))
	privKeyStr := marshalRSAPrivate(key)

	ioutil.WriteFile("./.root_keys/.priv", []byte(privKeyStr), 0600)
	ioutil.WriteFile("./.root_keys/.pub", []byte(pubKeyStr), 0600)

	return pubKeyStr, privKeyStr, nil
}

func encrypt(msg, publicKey string) (string, error) {
	parsed, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return "", err
	}
	// To get back to an *rsa.PublicKey, we need to first upgrade to the
	// ssh.CryptoPublicKey interface
	parsedCryptoKey := parsed.(ssh.CryptoPublicKey)

	// Then, we can call CryptoPublicKey() to get the actual crypto.PublicKey
	pubCrypto := parsedCryptoKey.CryptoPublicKey()

	// Finally, we can convert back to an *rsa.PublicKey
	pub := pubCrypto.(*rsa.PublicKey)

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		[]byte(msg),
		nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

func decrypt(data, priv string) (string, error) {
	data2, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode([]byte(priv))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, key, data2, nil)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

type SuperCrypt struct {
	pub  string
	priv string
}

func (s *SuperCrypt) Initialize() {

	_, err := ioutil.ReadFile("./.root_keys/.priv")
	if err != nil {
		if err == os.ErrNotExist {
			s.pub, s.priv, _ = generateKey()
		}
	}
}

func (s *SuperCrypt) Encrypt(msg string) (string, error) {
	return encrypt(msg, s.pub)
}

func (s *SuperCrypt) Decrypt(msg string) (string, error) {
	return decrypt(msg, s.priv)
}
