package account

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/holiman/uint256"
	"golang.org/x/crypto/ed25519"
	"log"
	"math/big"
	"strings"
)

type IAccount interface {
	New(args map[string]any) *Account
	NewFromKeys(id, ownerKey, activeKey string) *Account
}

type Account struct {
	ID       string
	Nickname string
	Groups   map[string]*Group
	ACL      map[string]*Permission
}

func (a Account) New(args map[string]any) *Account {
	if len(args) == 1 {
		return &Account{
			ID:     args["id"].(string),
			Groups: make(map[string]*Group),
			ACL:    make(map[string]*Permission),
		}
	}
	return &Account{}
}

func (a Account) NewFromKeys(id, ownerKey, activeKey string) *Account {
	acc := &Account{
		ID:     id,
		Groups: make(map[string]*Group),
		ACL:    make(map[string]*Permission),
	}
	acc.ACL["owner"] = &Permission{
		Name:      "owner",
		Threshold: 1,
		Items: []*GroupItem{
			{
				ID:        ownerKey,
				IsKeyPair: true,
				Weight:    1,
			},
		},
	}
	acc.ACL["active"] = &Permission{
		Name:      "active",
		Threshold: 1,
		Items: []*GroupItem{
			{
				ID:        activeKey,
				IsKeyPair: true,
				Weight:    1,
			},
		},
	}
	return acc
}

func NewAccountManager() IAccount {
	var A Account
	return A
}

func NewKeyPair(id string) (ed25519.PrivateKey, ed25519.PublicKey) {
	priv, pub := initializeKeyPair(id)
	return ed25519.PrivateKey(priv), ed25519.PublicKey(pub)
}

func NewAccountFromKeys(id, ownerKey, activeKey string) *Account {
	am := NewAccountManager()
	pub, priv := NewKeyPair(id)
	pb := priv
	pubb := pub
	a := am.NewFromKeys(id, hex.EncodeToString(pb), hex.EncodeToString(pubb))
	return a
}

func (a *Account) GetAddress() string {
	bb, _ := hex.DecodeString(a.ACL["active"].Items[0].ID)
	b := new(big.Int).SetBytes(bb)
	u, _ := uint256.FromBig(b)

	return strings.Replace(u.Hex(), "0x", "Qx00", -1)

}

func (a *Account) VerifyAddress(addr string) error {
	decoded, _ := uint256.FromHex(addr)
	hexBytes := decoded.Bytes()
	hexString := hex.EncodeToString(hexBytes)
	if hexString == a.ACL["active"].Items[0].ID {
		return nil
	}
	return errors.New("invalid address")
}

func (a *Account) Dump() {
	format := fmt.Sprintf("Account ID: %s \nAccount Address: %s", a.ID, a.GetAddress())
	log.Println(format)
	return
}

func (a *Account) String() string {
	return strings.Replace(a.GetAddress()[:38], "0x", "Qx00", -1)
}

func (a *Account) Bytes() []byte {
	return []byte("QBX" + "://" + a.GetAddress())
}
