package blockchain

import (
	"context"
	"github.com/google/uuid"
	"io/ioutil"
	"math/rand"
	"quantos/0.1.0/hash"
	"quantos/0.1.0/store"
	"quantos/0.1.0/version"
	"sync"
	"time"
)

type Blockchain struct {
	ID           string
	NetworkID    uint32
	SemVer       version.SemVer
	GenesisBlock []byte
	blocks       map[string]*Block
	blockHeaders map[string]*BlockHeader
	storage      store.Store
	ctx          context.Context
	manager      blockchainManager
}

func (b *Blockchain) GenerateID() {
	buf := make([]byte, 16)
	rand.NewSource(time.Now().UnixNano())
	rand.Read(buf)
	id := hash.NewHash(buf)
	uid, _ := uuid.FromBytes(id.Bytes())
	b.ID = uid.String()
}

func (b *Blockchain) SetVersion() {
	b.SemVer.Set(0, 1, 0)
}

func (b *Blockchain) CreateOrLoadGenesis() {
	bb, err := ioutil.ReadFile("./genesis.json")
	if err != nil {

		// create genesis block
		b.GenerateID()
		b.SetVersion()
		g := b.manager.Blocks().NewBlock()
		g.Head.BlockType = 0
		g.Head.BlockState = 0
		g.Head.ChainID = b.ID
		g.Head.Height = 1
		g.Head.ParentHash = hash.NewHash(nil).Bytes()
		g.Head.Timestamp = time.Now().UnixNano()
		g.Head.Version = 1
		g.Nonce = rand.Int()
		g.isFull.Store(false)
		g.payload = map[string][]byte{}
		g.Signatures = map[string][]byte{}
		b.manager.Blocks().blockObject = g
		b.manager.Blocks().WriteAsJson(g)

	} else {
		b.GenesisBlock = bb
	}
}

func NewBlockchain() *Blockchain {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	b := new(Blockchain)
	b.ctx = context.Background()
	b.CreateOrLoadGenesis()
	return b
}

func init() {
	NewBlockchain()
}
