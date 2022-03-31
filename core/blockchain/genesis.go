package blockchain

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/holiman/uint256"
	"github.com/quantosnetwork/dev-0.1.0/core/account"
	"github.com/quantosnetwork/dev-0.1.0/fs"
	"github.com/quantosnetwork/dev-0.1.0/hash"
	"github.com/quantosnetwork/dev-0.1.0/keygen/p2p"
	"github.com/quantosnetwork/dev-0.1.0/logger"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
	"github.com/wealdtech/go-merkletree"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io/ioutil"
	"lukechampine.com/blake3"
	"math/big"
	"strconv"
	"time"
)

var GenesisID = hex.EncodeToString(hash.NewBlake3Hash([]byte("genesis_kp")))

const LIVE_NETWORK = byte(0x0af)
const TEST_NETWORK = byte(0x0ba)
const CUSTOM_NET = byte(0x0f0)

type GenesisBlock struct {
	ctx               context.Context
	kctx              context.Context
	NetworkID         byte
	Account           *account.Account `json:"-"`
	Address           string
	CreationDate      int64
	Epoch             uint32
	GenesisValidators []*pb.Validator
	Block             *pb.Block
}

type GenesisContext struct {
	kv   map[string]any
	keys *GenKeys
}

type GenKeys struct {
	pub  p2p.P2PPublicKey
	priv p2p.P2PPrivateKey
}

type GKey string

var zeroHash = blake3.Sum256([]byte{0x00, 0x00, 0xff, 0x00, 0xff, 0x00, 0xff, 0xff})

func NewLiveGenesisBlock() *GenesisBlock {
	g := &GenesisBlock{
		ctx:          context.Background(),
		NetworkID:    LIVE_NETWORK,
		CreationDate: time.Now().UnixNano(),
		Account:      NewGenesisAccount(),
		Epoch:        0,
		Block:        new(pb.Block),
	}
	kctx := &GenesisContext{}
	nb := big.NewInt(0).SetBytes(zeroHash[:])
	b, _ := uint256.FromBig(nb)

	g.Block.Head = &pb.BlockHeader{
		BlockType:   pb.BlockType_GENESIS,
		Index:       0,
		Height:      1,
		ChainId:     uuid.New().String(),
		Version:     "1",
		Hash:        "",
		ParentHash:  b.Hex(),
		Timestamp:   timestamppb.New(time.Now()),
		BlockStates: []pb.BlockStates{pb.BlockStates_NEW, pb.BlockStates_PENDING_VALIDATION, pb.BlockStates_PENDING_SIGNATURES},
	}
	bb := make([][]byte, 1)
	bb[0] = []byte(g.Block.Head.ParentHash)
	merkle, _ := merkletree.New(bb)

	g.Block.Head.MerkleRoot = merkle.Root()
	key, keyp := g.NewKyberKeys()
	keys := &GenKeys{keyp, key}
	kctx.keys = keys
	g.Address = g.Account.String()
	g.Block.BlockId = g.Address

	g.Block.Head.Hash = NewBlockHash(g.Block)
	g.kctx = context.WithValue(g.ctx, GKey("genesis_context"), kctx)
	return g
}

func NewTestNetGenesisBlock(b *pb.Blockchain) {
	g := &GenesisBlock{
		ctx:       context.Background(),
		NetworkID: TEST_NETWORK,
		Account:   NewGenesisAccount(),
		Epoch:     0,
	}
	kctx := &GenesisContext{}

	key, keyp := g.NewKyberKeys()
	keys := &GenKeys{keyp, key}
	kctx.keys = keys
	g.kctx = context.WithValue(g.ctx, GKey("genesis_context"), kctx)
}

func NewCustomGenesisBlock(b *pb.Blockchain) {

}

func NewGenesisAccount() *account.Account {
	k := p2p.NewP2PKeys()
	pair := k.KeyPair()
	sk := pair.KeyPair().SK
	pk := pair.KeyPair().PK

	return account.NewAccountFromKeys(GenesisID, sk.String(), pk.String())

}

func (g *GenesisBlock) createGenesisHash() []byte {

	addr := g.Account.GetAddress()
	sig := g.Sign(addr)
	buf := make([][]byte, 34)
	buf[0] = []byte{g.NetworkID}
	buf[1] = g.Account.Bytes()
	buf[2], _ = g.Block.GetHead().Timestamp.AsTime().MarshalBinary()
	buf[3] = sig
	j := bytes.Join(buf, nil)
	return hash.NewKeyedBlake3Hash(j, sig)
}

func (g *GenesisBlock) Sign(msg string) []byte {
	keys := g.kctx.Value(GKey("genesis_context")).(GenesisContext)
	priv := keys.keys.priv
	s := priv.Sign(msg).GetSig().String()
	return []byte(s)
}

func (g *GenesisBlock) NewKyberKeys() (p2p.P2PPrivateKey, p2p.P2PPublicKey) {
	k := p2p.NewP2PKeys()
	pair := k.KeyPair()
	sk := pair.KeyPair().SK
	pk := pair.KeyPair().PK
	return sk, pk
}

func (g *GenesisBlock) SaveContext() error {

	file := &genCtxFile{}
	ctx := g.kctx.Value(GKey("genesis_context")).(*GenesisContext)
	b := ctx.keys
	file.Priv = b.priv.String()
	file.Pub = b.pub.String()
	j, _ := json.Marshal(file)
	enc := j
	err := ioutil.WriteFile(".genctx", enc, 0600)
	if err != nil {
		return err
	}
	return nil

}

type genCtxFile struct {
	encData []byte
	Pub     string
	Priv    string
}

func CheckIfGenesisFileExists(networkID byte) bool {
	qfs := fs.NewFileSystem()

	netstr := strconv.Itoa(int(networkID))
	exists, err := qfs.Exists("./genesis-" + netstr + ".json")
	if err != nil {
		return false
	}
	if !exists {
		return false
	}
	return true
}

func WriteGenesisBlock(networkID byte, genesis *GenesisBlock) error {

	g, err := json.Marshal(genesis)
	if err != nil {
		logger.Logger.Sugar().Fatal("an error happened writing genesis block", err)
		return err
	}
	netstr := strconv.Itoa(int(networkID))
	err = ioutil.WriteFile("./genesis-"+netstr+".json", g, 0600)
	if err != nil {
		return err
	}
	return genesis.SaveContext()

}
