package core

import (
	"github.com/holiman/uint256"
	"github.com/quantosnetwork/dev-0.1.0/common"
	"github.com/quantosnetwork/dev-0.1.0/keygen/p2p"
)

type EncodingUtils interface {
	GetEncoder() common.Encoder
	ToUint256()
}

type Uint256Util interface {
	FromBytes([]byte) *uint256.Int
	FromFloat64(f float64) *uint256.Int
	FromInt64(i int64) *uint256.Int
	FromString(s string) *uint256.Int
	GetAddressFromType(from string, data []byte) string
	SafeMath(operation string, subject map[string]any, result *uint256.Int) error
}

type address []byte

type Address interface {
	GenerateNew() *Address
	String() string
	QBX() string
	QBT() string
	Keys() *p2p.Keys
	Decode() *address
	Encode() []byte
	Check(a *Address) bool
	FromAccount(accId string) *Address
	Derive() *Address
}
