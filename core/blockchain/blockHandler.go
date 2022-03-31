package blockchain

import (
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
)

/*
	@description
	BlockHandler is the function managing all non-immutable block modifications via the VBlock proxy
*/

type BlockV1 pb.Block
type VBlock pb.Block

type BlockHandler interface {
	Modify(*BlockV1, *VBlock)
	Parse(*BlockV1, *VBlock)
	Verify(*BlockV1, *VBlock)
	Apply(original *BlockV1, blockEdits *BlockV1)
}

type BlockHandlerFunc func(*BlockV1, *VBlock)

func (b BlockHandlerFunc) Parse(original *BlockV1, proxy *VBlock) {
	b(original, proxy)
}

func (b BlockHandlerFunc) Modify(original *BlockV1, proxy *VBlock) {
	b(original, proxy)
}

func (b BlockHandlerFunc) Verify(original *BlockV1, proxy *VBlock) {
	b(original, proxy)
}

func (b BlockHandlerFunc) Apply(original *BlockV1, blockEdits *BlockV1) {
	b(original, &VBlock{})
}
