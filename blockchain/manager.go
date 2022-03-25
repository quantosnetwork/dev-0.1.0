package blockchain

import (
	"github.com/hashicorp/golang-lru"
	"github.com/looplab/fsm"
	"github.com/quantosnetwork/v0.1.0-dev/blockchain/block"
	"github.com/quantosnetwork/v0.1.0-dev/config"
	"github.com/quantosnetwork/v0.1.0-dev/tx"
)

const MaxTxPerBlock = 100
const MaxContractPerBlock = 30

type Manager interface {
	CreateNewBlock(chainID string, creator string, block *block.BlockV1) (block.Block, error)
	CreateGenesis() error
	GetBlockHeaders() []*block.BlockHeader
	GetReceipts() []interface{}
	GetAllBlocks() []block.Block
	GetBlockFromHeight(height uint64)
	GetGenesisBlock() block.Block
	GetBlockByHeight(height uint64) (block.Block, error)
	GetBlockByIndex(index uint64) (block.Block, error)
	GetBlockByTxID(txID string) (block.Block, error)
	GetCurrentState()
	GetLastBlock() block.Block
	SignBlock(block *block.BlockV1)
	ValidateBlockchain(chainID string) bool
	GetCoinbase(currencyID uint)
	Accounts() // returns account manager
	Blocks() *block.VBlock
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

func (b blockchainManager) CreateNewBlock(chainID string, creator string, block *block.BlockV1) (block.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) CreateGenesis() error {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockHeaders() []*block.BlockHeader {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetReceipts() []interface{} {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetAllBlocks() []block.Block {
	bs := make([]block.Block, 1)
	return bs
}

func (b blockchainManager) GetBlockFromHeight(height uint64) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetGenesisBlock() block.Block {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockByHeight(height uint64) (block.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockByIndex(index uint64) (block.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetBlockByTxID(txID string) (block.Block, error) {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetCurrentState() {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) GetLastBlock() block.Block {
	//TODO implement me
	panic("implement me")
}

func (b blockchainManager) SignBlock(block *block.BlockV1) {
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
func (b blockchainManager) Blocks() *block.VBlock {
	return block.GetBlockProxy()
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
