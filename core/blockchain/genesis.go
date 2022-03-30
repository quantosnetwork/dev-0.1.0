package blockchain

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/quantosnetwork/dev-0.1.0/core/account"
	"github.com/quantosnetwork/dev-0.1.0/crypto/anon"
	"github.com/quantosnetwork/dev-0.1.0/hash"
	"github.com/quantosnetwork/dev-0.1.0/keygen/p2p"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
	"io/ioutil"
)

var GenesisID = hex.EncodeToString(hash.NewBlake3Hash([]byte("genesis_kp")))

const LIVE_NETWORK = byte(0x0af)
const TEST_NETWORK = byte(0x0ba)
const CUSTOM_NET = byte(0x0f0)

type GenesisBlock struct {
	ctx               context.Context
	kctx              context.Context
	NetworkID         byte
	Account           *account.Account
	CreationDate      int64
	Epoch             uint32
	GenesisValidators []*pb.Validator
	block             *pb.Block
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

func NewLiveGenesisBlock(b *pb.Blockchain) {
	g := &GenesisBlock{
		ctx:       context.Background(),
		NetworkID: LIVE_NETWORK,
		Account:   NewGenesisAccount(),
		Epoch:     0,
	}
	kctx := &GenesisContext{}

	key, keyp := g.NewKyberKeys()
	keys := &GenKeys{keyp, key}
	kctx.keys = keys
	g.kctx = context.WithValue(g.ctx, GKey("genesis_context"), kctx)
}

func NewTestNetGenesisBlock(b *pb.Blockchain) {

}

func NewCustomGenesisBlock(b *pb.Blockchain) {

}

func NewGenesisAccount() *account.Account {
	priv, pub := account.NewKeyPair(GenesisID)
	sks := hex.EncodeToString(priv)
	pks := hex.EncodeToString(pub)
	return account.NewAccountFromKeys(GenesisID, sks, pks)

}

func (g *GenesisBlock) createGenesisHash() []byte {

	addr := g.Account.GetAddress()
	sig := g.Sign(addr)
	buf := make([][]byte, 34)
	buf[0] = []byte{g.NetworkID}
	buf[1] = g.Account.Bytes()
	buf[2], _ = g.block.GetHead().Timestamp.AsTime().MarshalBinary()
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

	nonce := make([]byte, 8)
	f := func(pktNum uint64, in []byte) ([]byte, []byte, bool) {
		binary.BigEndian.PutUint64(nonce, pktNum)
		var buf [32]byte
		ctx := g.kctx.Value(GKey("genesis_context")).(GenesisContext)
		b, _ := ctx.keys.priv.Bytes()
		copy(buf[:], b)
		encoded, err := anon.EncodeNECS(&buf, nonce, in)
		if err != nil {
			return nil, nil, false
		}
		decoded, err := anon.NECSDecode(&buf, nonce, encoded)

		if err != nil {

			return nil, nil, false
		}

		return encoded, decoded, bytes.Compare(decoded, in) == 0
	}

	file := &genCtxFile{}
	ctx := g.kctx.Value(GKey("genesis_context")).(GenesisContext)
	b := ctx.keys
	file.Priv = b.priv.String()
	file.Pub = b.pub.String()
	j, _ := json.Marshal(file)
	enc, _, _ := f(0, j)
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
