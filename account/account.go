package account

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cloudflare/circl/sign/ed25519"
	"github.com/quantosnetwork/v0.1.0-dev/common"
	"log"
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

type GroupItem struct {
	ID         string //keypair id
	Permission string
	IsKeyPair  bool
	Weight     int
}

type Group struct {
	Name  string
	Items []*GroupItem
}

type Permission struct {
	Name      string
	Groups    []string
	Items     []*GroupItem
	Threshold int
}

func NewAccountManager() IAccount {
	var A Account
	return A
}

func NewKeyPair(id string) (ed25519.PrivateKey, ed25519.PublicKey) {
	priv, pub := initializeKeyPair(id)
	return priv, pub
}

func NewAccountFromKeys(id, ownerKey, activeKey string) *Account {
	am := NewAccountManager()
	pub, priv := NewKeyPair(id)
	pb, _ := priv.MarshalBinary()
	pubb, _ := pub.MarshalBinary()
	a := am.NewFromKeys(id, hex.EncodeToString(pb), hex.EncodeToString(pubb))
	return a
}

func (a *Account) GetAddress() string {

	return common.EncodeBase58([]byte(a.ACL["active"].Items[0].ID))

}

func (a *Account) VerifyAddress(addr string) error {
	decoded := common.DecodeBase58(addr)
	hexString := hex.EncodeToString(decoded)
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
	return "QBX" + "://" + a.GetAddress()
}

func (a *Account) Bytes() []byte {
	return []byte("QBX" + "://" + a.GetAddress())
}
