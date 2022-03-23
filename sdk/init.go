package sdk

import (
	"github.com/quantosnetwork/v0.1.0-dev/account"
)

type QuantosSDK interface {
	Accounts() AccountManager
}

var Q QuantosSDK
var LoadedAccount *account.Account

func InitializeSDK() {
	LoadedAccount = Q.Accounts().GetLoadedAccount()
}
