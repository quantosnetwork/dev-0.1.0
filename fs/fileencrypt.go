package fs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

type FileEncryptor struct {
	buf  bytes.Buffer
	priv *rsa.PrivateKey
	pub  rsa.PublicKey
}

func (fe *FileEncryptor) NewFileKeys() {

	keys, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := keys.PublicKey
	fe.priv = keys
	fe.pub = pub

}

func (fe *FileEncryptor) Encrypt(data []byte) ([]byte, error) {
	fe.buf.Reset()
	enc, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		&fe.pub,
		data,
		nil,
	)

	fe.buf.Write(enc)
	return fe.buf.Bytes(), err
}

func (fe *FileEncryptor) Decrypt() ([]byte, error) {
	encB := fe.buf.Bytes()
	dec, err := fe.priv.Decrypt(nil, encB, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func NewFileEncryptorWithNewKeys() *FileEncryptor {
	f := new(FileEncryptor)
	f.NewFileKeys()
	return f
}

func NewFileEncryptorWithLoadedKeys(lPub interface{}, lPriv interface{}) *FileEncryptor {
	priv := lPriv.(*rsa.PrivateKey)
	pub := lPub.(rsa.PublicKey)
	return &FileEncryptor{
		priv: priv,
		pub:  pub,
	}
}
