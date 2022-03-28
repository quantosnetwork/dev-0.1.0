package blockchain

import (
	"context"
	"errors"
	"github.com/quantosnetwork/dev-0.1.0/core/account"
	"sync"
)

// BlockProxy protects the existing block by filtering requests and avoiding unwanted tempering
type BlockProxy interface {
	Initialize(chainID string, callerAccount account.Account) error
	IsAllowed(aclType int, callerAccount account.Account) bool
	CreateBlankBLock() *BlockV1
	CopyBlock(b *BlockV1) *BlockV1
	FinalizeBlock(b *BlockV1) error
	GetRawBlock(height uint32) (*BlockV1, error)
	AddBlockToValidationQueue(b *BlockV1)
	ValidateAndSign(b *BlockV1) error
}

type BlockProxyQueries uint32

const (
	BPQ BlockProxyQueries = iota
	INIT
	CREATEBLANK
	COPY
	FINALIZE
	GETRAW
	ADDTOVALQUEUE
	VALANDSIGN
)

type Proxy struct {
	isInitialized chan bool
	ctx           context.Context
	chainID       string
	caller        account.Account
	guardFunc     GuardFunc
	mu            sync.RWMutex
	BlockProxy
}

type ProxyContextKey string

type GuardFunc = func(caller account.Account, query BlockProxyQueries) error
type RunFunc = func(g GuardFunc, query BlockProxyQueries, data ...any) (any, error)

func (bp *Proxy) Initialize(chainID string, callerAccount account.Account) error {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	acl := callerAccount.ACL["proxy"]
	if acl.Name != "allow_proxy_requests" {
		return errors.New("not allowed to request block proxy")
	}

	return nil
}
