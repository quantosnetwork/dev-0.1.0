package sdk

import (
	"github.com/quantosnetwork/dev-0.1.0/core/account"
)

type QuantosSDK interface {
	Accounts() AccountManager
}

var Q QuantosSDK
var LoadedAccount *account.Account

func InitializeSDK() {
	LoadedAccount = Q.Accounts().GetLoadedAccount()
}
