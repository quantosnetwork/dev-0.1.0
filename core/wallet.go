package core

import (
	"github.com/holiman/uint256"
)

type Wallet interface {
	Address() *Address
	Lock() bool
	Unlock(key []byte) bool
	Backup() (string, error)
	Restore(backupKey []byte, address string, mnemonic string) (bool, error)
	Withdraw()
	SendTX(from Address, to Address, amt uint256.Int) (txId string, success bool, err error)
	Synchronize()
	Balance() *uint256.Int
}
