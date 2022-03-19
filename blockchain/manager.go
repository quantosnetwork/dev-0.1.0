package blockchain

import (
	"github.com/hashicorp/golang-lru"
	"github.com/looplab/fsm"
	"quantos/0.1.0/config"
	"quantos/0.1.0/tx"
)

const MaxTxPerBlock = 100
const MaxContractPerBlock = 30

type Manager interface {
	CreateNewBlock(chainID string, creator string, block *BlockV1) (Block, error)
	CreateGenesis() error
	GetBlockHeaders() []*BlockHeader
	GetReceipts() []interface{}
	GetAllBlocks() []Block
	GetBlockFromHeight(height uint64)
	GetGenesisBlock() Block
	GetBlockByHeight(height uint64) (Block, error)
	GetBlockByIndex(index uint64) (Block, error)
	GetBlockByTxID(txID string) (Block, error)
	GetCurrentState()
	GetLastBlock() Block
	SignBlock(block *BlockV1)
	ValidateBlockchain(chainID string) bool
	GetCoinbase(currencyID uint)
	Accounts() // returns account manager
	Blocks() *VBlock
	Consensus()       // returns consensus manager
	States()          // returns states manager
	Tx() tx.TxManager // returns tx manager
	Bank()            // returns Coinbase manager (Bank functions)
	BlockQueue()      // returns block queue
	TxQueue()         // returns tx queue
	Config() *config.ChainConfig
	FSM() *fsm.FSM
	Cache() *lru.ARCCache
	Stats()   // returns chain info / stats / metrics
	Bridger() // returns bridge interface for multichains
	Indexer() // returns indexers interface to speed up block finder
}

type blockchainManager struct{}

func (b blockchainManager) CreateNewBlock(chainID string, creator string, block *BlockV1) (Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) CreateGenesis() error {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockHeaders() []*BlockHeader {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetReceipts() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetAllBlocks() []Block {
	bs := make([]Block, 1)
	return bs
}

func (b blockchainManager) GetBlockFromHeight(height uint64) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetGenesisBlock() Block {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockByHeight(height uint64) (Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockByIndex(index uint64) (Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockByTxID(txID string) (Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetCurrentState() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetLastBlock() Block {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) SignBlock(block *BlockV1) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) ValidateBlockchain(chainID string) bool {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetCoinbase(currencyID uint) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Accounts() {
	//TODO implement me
	panic("implement me")
}

// Blocks blockchainManager.Blocks is an alias to blockchainManager.GetAllBlocks will return all Block interfaces
func (b blockchainManager) Blocks() *VBlock {
	return GetBlockProxy()
}

func (b blockchainManager) Consensus() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) States() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Tx() tx.TxManager {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Bank() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) BlockQueue() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) TxQueue() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Config() *config.ChainConfig {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) FSM() *fsm.FSM {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Cache() *lru.ARCCache {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Stats() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Bridger() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) Indexer() {
	//TODO implement me
	panic("implement me")
}

func NewBlockchainManager() Manager {
	return &blockchainManager{}
}
