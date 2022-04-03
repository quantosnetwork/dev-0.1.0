package fs

import (
	"encoding/json"
	"github.com/quantosnetwork/dev-0.1.0/crypto/super"
)

var Crypt *super.SuperCrypt

func init() {
	s := &super.SuperCrypt{}
	s.Initialize()
	Crypt = s

}

func (read *Reader) Decrypt(data string) (string, error) {
	return Crypt.Decrypt(data)
}

func (read *Reader) DecryptToStruct(data string) (any, error) {
	m, err := read.Decrypt(data)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal([]byte(m), &read.Options.Struc)
	return read.Options.Struc, nil
}

func (w *Writer) Encrypt(data string) (string, error) {
	return Crypt.Encrypt(data)
}

func (w *Writer) EncryptFromStruct(data any) ([]byte, error) {
	d, _ := json.Marshal(data)
	m, err := w.Encrypt(string(d))
	if err != nil {
		return nil, err
	}
	mb := make([]byte, len([]byte(m)))
	copy(mb, []byte(m))

	return mb, nil
}
