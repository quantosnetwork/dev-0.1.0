package block

/*
	@description
	BlockHandler is the function managing all non-immutable block modifications via the VBlock proxy
*/

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
	b(original, &VBlock{BlockImage: blockEdits})
}
