package sdk

import (
	"context"
	"encoding/hex"
	"github.com/quantosnetwork/v0.1.0-dev/account"
)

type AccountManager interface {
	CheckIfLoadedAccount(id string) *accountManager
	GetLoadedAccount() *account.Account
	GetLoadedKeys() *account.Keys
}

type accountManager struct {
	Ctx  context.Context
	Keys *account.LoadedKeys
	ID   string
	AccountManager
}

func (a accountManager) CheckIfLoadedAccount(id string) *accountManager {
	keys := account.Keys{}
	lk := keys.GetLoadedKeys(id)

	return &accountManager{
		ID:   id,
		Keys: lk,
	}
}

func (a *accountManager) GetLoadedAccount() *account.Account {
	pb, _ := a.Keys.Pub.MarshalBinary()
	sk, _ := a.Keys.Priv.MarshalBinary()
	acct := account.NewAccountFromKeys(a.ID, hex.EncodeToString(sk), hex.EncodeToString(pb))
	return acct
}

func NewAccountManager(ctx context.Context, accID string) *accountManager {
	var am accountManager
	accCtx := context.WithValue(ctx, []byte("aID"), accID)
	am.Ctx = accCtx
	a := am.CheckIfLoadedAccount(accID)
	return a
}
