package sdk

import (
	"context"
	"encoding/hex"
	account2 "github.com/quantosnetwork/dev-0.1.0/core/account"
)

type AccountManager interface {
	CheckIfLoadedAccount(id string) *accountManager
	GetLoadedAccount() *account2.Account
	GetLoadedKeys() *account2.Keys
}

type accountManager struct {
	Ctx  context.Context
	Keys *account2.LoadedKeys
	ID   string
	AccountManager
}

func (a accountManager) CheckIfLoadedAccount(id string) *accountManager {
	keys := account2.Keys{}
	lk := keys.GetLoadedKeys(id)

	return &accountManager{
		ID:   id,
		Keys: lk,
	}
}

func (a *accountManager) GetLoadedAccount() *account2.Account {
	pb, _ := a.Keys.Pub.MarshalBinary()
	sk, _ := a.Keys.Priv.MarshalBinary()
	acct := account2.NewAccountFromKeys(a.ID, hex.EncodeToString(sk), hex.EncodeToString(pb))
	return acct
}

func NewAccountManager(ctx context.Context, accID string) *accountManager {
	var am accountManager
	accCtx := context.WithValue(ctx, []byte("aID"), accID)
	am.Ctx = accCtx
	a := am.CheckIfLoadedAccount(accID)
	return a
}
