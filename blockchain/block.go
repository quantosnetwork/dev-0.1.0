package blockchain

import (
	"github.com/barkimedes/go-deepcopy"
	"go.uber.org/atomic"
	"quantos/0.1.0/hash"
	"quantos/0.1.0/tx"
)

type Block interface {
	Type() uint64
	Height() uint64
	Index() uint64
	ParentHash() *hash.Hash
	Version() int32
	Hash() *hash.Hash
	Payload() map[string][]byte
	BlockTime() int64
	Transactions() map[string]tx.Transaction
	HasTX(txId string) bool
	Number() int
	Finalized() bool
	Size() int
	GetRaw() []byte
	Header() *BlockHeader
	ChainID() string
	Proxy() *VBlock
}

type BlockHeader struct {
	BlockType         uint64
	Index             uint64
	Height            uint64
	ChainID           string
	Version           int32
	Hash              *hash.Hash
	ParentHash        *hash.Hash
	MerkleRoot        *hash.Hash
	TxMerkleRoot      *hash.Hash
	ReceiptMerkleRoot *hash.Hash
	Timestamp         int64
	Number            int
	Size              int
	NumTx             int
	TxIds             map[string]bool
	CreationTx        string
	BlockState        uint32
}

type BlockV1 struct {
	Head           *BlockHeader
	payload        map[string][]byte
	OpenedTxSlots  [MaxTxPerBlock]map[string]*tx.Transaction
	Signatures     map[string][]byte
	ContractsSlots [MaxContractPerBlock]map[string]interface{}
	Nonce          int
	Validators     map[string]bool
	isFull         atomic.Bool
}

func (b BlockV1) Type() uint64 {
	return b.Head.BlockType
}

func (b BlockV1) Height() uint64 {
	return b.Head.Height
}

func (b BlockV1) Index() uint64 {
	return b.Head.Index
}

func (b BlockV1) ParentHash() *hash.Hash {
	return b.Head.ReceiptMerkleRoot
}

func (b BlockV1) Version() int32 {
	return b.Head.Version
}

func (b BlockV1) Hash() *hash.Hash {
	return b.Head.Hash
}

func (b BlockV1) Payload() map[string][]byte {
	return b.payload
}

func (b BlockV1) BlockTime() int64 {
	return b.Head.Timestamp
}

func (b BlockV1) Transactions() map[string]tx.Transaction {
	//TODO implement me
	panic("implement me")
}

func (b BlockV1) HasTX(txId string) bool {
	//TODO implement me
	panic("implement me")
}

func (b BlockV1) Number() int {
	return b.Head.Number
}

func (b BlockV1) Finalized() bool {
	return b.isFull.Load()
}

func (b BlockV1) Size() int {
	return b.Head.Size
}

func (b BlockV1) GetRaw() []byte {
	//TODO implement me
	panic("implement me")
}

func (b BlockV1) Header() *BlockHeader {
	return b.Head
}

func (b BlockV1) ChainID() string {
	return b.Head.ChainID
}

func (b BlockV1) Proxy() *VBlock {
	return &VBlock{blockObject: &b}
}

// VBlock is the block Proxy
type VBlock struct {
	blockObject *BlockV1
	BlockImage  *BlockV1
}

func (vb *VBlock) NewBlock() Block {
	b := &BlockV1{}
	bh := &BlockHeader{}
	b.Head = bh
	vb.blockObject = b
	return b
}

// duplicates the original block (perfectly deepcopy) so we ca work on the copy and revert in case of an error
// without disrupting the chain
func (vb *VBlock) duplicateOriginalObject() {

	blockCopy, err := deepcopy.Anything(vb.blockObject)
	if err != nil {
		panic(err)
	}
	vb.BlockImage = blockCopy.(*BlockV1)

}

// GetBlockProxy will get the block interface to perform actions on the real block
func GetBlockProxy(blockIndex int64) {

}
