package config

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
	"github.com/quantosnetwork/v0.1.0-dev/keygen/p2p"
	"github.com/quantosnetwork/v0.1.0-dev/serializer"
	"github.com/quantosnetwork/v0.1.0-dev/uptime"
	"github.com/quantosnetwork/v0.1.0-dev/version"
	"go.uber.org/atomic"
	"io/ioutil"
	"time"
)

const DEFAULTHOST = ":55655"

type NodeConfig struct {
	ID            *ID
	Version       version.Version
	ListenAddress string
	Logger        hclog.Logger
	Serializer    *serializer.SerializableItem
	Worker        interface{}
	NodeState     atomic.Int64
	Uptime        uptime.UptimeManager
	Keys          p2p.P2PKeys
	seed          []byte

	QuitCh chan struct{}
}

type NodeContext struct {
	ParentCtx context.Context
	Config    *NodeConfig
	SessionID uuid.UUID
}

func NewNodeConfig(hostAddr string) *NodeConfig {
	v := &version.SemVer{}
	v.Set(0, 1, 0)
	err := ioutil.WriteFile(".version", v.Hash(), 0644)
	if err != nil {
		panic(err)
	}

	nc := &NodeConfig{}
	nc.Version = v
	nc.ListenAddress = hostAddr
	nc.Logger = hclog.New(hclog.DefaultOptions)
	nc.Serializer = new(serializer.SerializableItem)
	nc.NodeState.Store(0)
	nc.Uptime = uptime.Manager()
	nc.ID = nc.NewID()
	return nc
}

type ID struct {
	r []byte
	s []byte
	d string
	g []byte
}

func (nc *NodeConfig) NewID() *ID {

	nc.Keys = p2p.NewP2PKeys()
	nc.seed = IDSeed()
	pubString, _point, _scalar := nc.Keys.KeyPair().DerivePubKey(hex.EncodeToString(nc.seed))
	buf := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(pubString))))
	base64.StdEncoding.Encode(buf, []byte(pubString))
	pub := buf
	r, _ := _point.MarshalBinary()
	s, _ := _scalar.MarshalBinary()
	id := &ID{
		r: r,
		s: s,
		d: string(pub),
		g: nc.seed,
	}

	return id
}

func (id *ID) String() string {
	mid, _ := json.Marshal(id)
	hasher := sha256.New()
	hasher.Write(mid)
	s := hasher.Sum(nil)
	return hex.EncodeToString(s)
}

func (nc *NodeConfig) ValidateID(id *ID) error {

	return nil
}

func IDSeed() []byte {
	buf := make([]byte, 64)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	return buf[:]

}

type OnDiskConfig struct {
	NodeID struct {
		Raw []byte
	}
	Version   string `json:"node_version"`
	State     int64  `json:"last_node_state"`
	SessionID string `json:"last_session_id"`
	Keys      struct {
		SK string
		PK string
	}
	SavedOn int64
}

func WriteConfigItem(ctx context.Context) {

	var nctx *NodeContext
	nctx = ctx.Value("nctx").(*NodeContext)
	tw := new(OnDiskConfig)
	buf := make([][]byte, 3)
	buf[0] = nctx.Config.ID.r
	buf[1] = nctx.Config.ID.s
	buf[2] = nctx.Config.ID.g
	tw.NodeID.Raw = bytes.Join(buf, []byte("_"))
	tw.Version = nctx.Config.Version.String()
	tw.State = nctx.Config.NodeState.Load()
	tw.Keys.SK = nctx.Config.Keys.KeyPair().SK.String()
	tw.Keys.PK = nctx.Config.Keys.KeyPair().PK.String()
	tw.SavedOn = time.Now().UnixNano()
	twb, _ := json.Marshal(tw)
	ioutil.WriteFile("sessions/"+nctx.SessionID.String()+".json", twb, 0600)

}

func NewNodeContext(ctx context.Context) context.Context {
	nctx := &NodeContext{}
	nctx.Config = NewNodeConfig(DEFAULTHOST)
	nctx.ParentCtx = ctx
	sess, err := uuid.FromBytes(nctx.Config.seed[0:16])
	if err != nil {
		panic(err)
	}
	nctx.SessionID = sess
	c := context.WithValue(nctx.ParentCtx, "nctx", nctx)
	WriteConfigItem(c)

	return c
}
