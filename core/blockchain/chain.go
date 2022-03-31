package blockchain

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/holiman/uint256"
	"github.com/looplab/fsm"
	"github.com/quantosnetwork/dev-0.1.0/core/trie"
	"github.com/quantosnetwork/dev-0.1.0/hash"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
	"github.com/quantosnetwork/dev-0.1.0/store"
	"github.com/quantosnetwork/dev-0.1.0/version"
	"io/ioutil"
	"lukechampine.com/frand"
	"math/big"
	"math/rand"
	"strings"
	"sync"
	"time"
)

const STAKE_COST_PER_BLOCK = 10

type Blockchain struct {
	SemVer       version.SemVer
	GenesisBlock []byte
	blocks       map[string]*pb.Block
	blockHeaders map[string]*pb.BlockHeader
	storage      store.Store
	ctx          context.Context
	ChainStates  map[string]*fsm.FSM
	BlockVals    *BlockValidation
	*pb.Blockchain
}

func (b *Blockchain) GenerateID() {
	buf := make([]byte, 16)
	rand.NewSource(time.Now().UnixNano())
	rand.Read(buf)
	id := hash.NewHash(buf)
	uid, _ := uuid.FromBytes(id.Bytes())
	b.ChainId = uid.String()
}

func (b *Blockchain) SetVersion() {
	b.SemVer.Set(0, 1, 0)
}

func (b *Blockchain) CreateOrLoadGenesis() {

}

func (b *Blockchain) InitializeStateMachines(machineNames []string, events map[string][]fsm.EventDesc,
	cbs map[string]map[string]fsm.Callback) {
	for _, m := range machineNames {
		evts := events[m]
		callbacks := cbs[m]
		b.ChainStates[m] = fsm.NewFSM(
			"idle", evts, callbacks)
	}
}

func NewBlockchain(networkId byte) *pb.Blockchain {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	b := new(pb.Blockchain)

	genesis := NewLiveGenesisBlock()
	b.ChainId = genesis.Block.GetHead().GetChainId()
	b.BlockchainHead = genesis.Block
	b.GenesisHash = genesis.Block.GetHead().GetHash()
	b.ChainId = b.GetBlockchainHead().GetHead().ChainId
	b.Version = "1"
	b.NumBlocks = 1

	b.Blocks = append(b.Blocks, genesis.Block)

	bb := big.NewInt(0).SetBytes(b.BlockchainHead.GetHead().MerkleRoot)
	bmk, _ := uint256.FromBig(bb)
	b.MerkleRoot = bmk.Hex()
	bc, _ := json.Marshal(b)
	ioutil.WriteFile("./chain.json", bc, 0600)
	return b
}

func init() {
	NewBlockchain(LIVE_NETWORK)
}

func (b *Blockchain) GenerateNewBlock(validator *pb.Validator) ([]*pb.Block, *pb.Block, error) {
	if err := b.ValidateBlockchain(); err != nil {
		val := validator.GetNode()
		stake := val.GetStake()
		stake -= STAKE_COST_PER_BLOCK
		return b.Blockchain.Blocks, b.BlockchainHead, err
	}

	nb := &pb.Block{}
	nbh := &pb.BlockHeader{}
	nbh.BlockStates = append(nbh.BlockStates, pb.BlockStates_NEW,
		pb.BlockStates_PENDING_VALIDATION, pb.BlockStates_PENDING_SIGNATURES)
	nbh.Version = b.Version
	nbh.ChainId = b.GetChainId()
	nbh.ParentHash = b.GetBlockchainHead().GetHead().GetHash()
	nbh.Height = b.GetBlockchainHead().GetHead().GetHeight() + 1
	nbh.Index = b.GetBlockchainHead().GetHead().GetIndex() + 1

	nbh.MerkleRoot = b.calculateMerkleTree()
	nbh.Number = int32(b._genRandomBlockNumber())
	nbh.Size = int64(0)
	nbh.NumTx = int64(0)
	nb.Head = nbh
	nb.Nonce = b._genRandomBlockNumber()
	nb.ValidatorAddr = validator.String()
	bid := NewBlockHash(nb)
	nb.BlockId = bid
	nbh.Hash = NewBlockHash(nb)

	if err := b.ValidateBlockCandidate(nb); err != nil {
		stake := validator.GetNode().GetStake()
		stake -= STAKE_COST_PER_BLOCK
		return b.Blocks, nb, nil
	} else {
		b.blocks[bid] = nb
	}

	b.blockHeaders[nbh.Hash] = nbh

	return b.Blocks, nb, nil

}

func NewBlockHash(nb *pb.Block) string {
	toHash := []byte(nb.GetHead().GetTimestamp().String() + nb.GetHead().ParentHash + nb.
		GetValidatorAddr())
	h := hash.NewBlake3Hash(toHash)
	b := big.NewInt(0).SetBytes(h)
	u, _ := uint256.FromBig(b)
	return u.Hex()
}

type BlockValidation struct {
	Invalid map[string]error
	Valid   map[string]bool
}

func (b *Blockchain) ValidateBlock(nb *pb.Block) bool {
	prevHash := b.GetBlockchainHead().GetHead().GetHash()
	ts := nb.GetHead().GetTimestamp().String()
	valaddr := nb.GetValidatorAddr()

	err := VerifyHash(nb.GetHead().GetHash(), prevHash, ts, valaddr)
	if err != nil {
		b.BlockVals.Invalid[nb.GetBlockId()] = err
	}

	prevIndex := b.GetBlockchainHead().GetHead().GetIndex()
	if nb.GetHead().GetIndex() != prevIndex+1 {
		b.BlockVals.Invalid[nb.GetBlockId()] = errors.New("block index is invalid")
	}
	prevHeight := b.GetBlockchainHead().GetHead().GetHeight()
	if nb.GetHead().GetIndex() != prevHeight+1 {
		b.BlockVals.Invalid[nb.GetBlockId()] = errors.New("block height is invalid")
	}

	prevTime := b.GetBlockchainHead().GetHead().GetTimestamp().AsTime()
	current := nb.GetHead().GetTimestamp().AsTime()
	if current.After(prevTime) {
		b.BlockVals.Valid[nb.GetBlockId()] = true
		return true
	} else {
		b.BlockVals.Invalid[nb.GetBlockId()] = errors.New("block time is invalid")
	}
	return false
}

func (b *Blockchain) ValidateBlockchain() error {
	if len(b.blocks) <= 1 {
		return nil
	}

	for _, bl := range b.blocks {
		valid := b.ValidateBlock(bl)
		if !valid {
			return errors.New("blockchain is invalid")
		}
	}
	return nil
}

func VerifyHash(currentHash string, prevHash string, ts string, valaddr string) error {
	toCompare := []byte(ts + prevHash + valaddr)
	tch := hash.NewBlake3Hash(toCompare)
	tchs := hex.EncodeToString(tch)
	if strings.EqualFold(currentHash, tchs) {
		return nil
	}
	return errors.New("invalid hash")
}

func (b *Blockchain) calculateMerkleTree() []byte {

	merkleRootData := b.Blocks
	blockmerkledata := make([][]byte, len(merkleRootData))

	for i, bl := range merkleRootData {
		blockmerkledata[i] = []byte(bl.GetHead().GetHash())
	}

	mr := trie.NewMerkleTree(b.GetChainId(), string(b.GetVersion()), blockmerkledata...)
	return mr.Tree.Root()

}

func (b *Blockchain) _genRandomBlockNumber() uint32 {
	buf := make([]byte, 16)
	frand.Read(buf)
	return binary.LittleEndian.Uint32(buf)
}

func (b *Blockchain) ValidateBlockCandidate(newBlock *pb.Block) error {
	valid := b.ValidateBlock(newBlock)
	if !valid {
		return errors.New("validator failed to validate block")
	}
	return nil
}

func (b *Blockchain) TriggerStateEvent(stateType string, evt string, args ...any) error {
	err := b.ChainStates[stateType].Event(evt, args)
	if err != nil {
		return err
	}
	return nil
}
