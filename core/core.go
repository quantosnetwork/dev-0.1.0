package core

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/quantosnetwork/dev-0.1.0/core/blockchain"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
	"github.com/quantosnetwork/dev-0.1.0/version"
	"strconv"
)

/*
	Core interface to manage everything
*/

type Chain interface {
	Name() string
	Genesis() GenesisBlock
	Params() map[string]any
	BootNodes() []string
}

type ChainGetter interface {
	Import(chain string) (pb.Blockchain, error)
	ImportFromChainId(chainId string) (pb.Blockchain, error)
	ImportFromFilePath(filepath string) (pb.Blockchain, error)
	ImportFromDatabase(db *badger.DB) (pb.Blockchain, error)
}

var Version version.SemVer

var Networks = map[uint32]string{
	175:  "LivePublic",
	177:  "LivePrivate",
	199:  "Testnet",
	0:    "Local",
	999:  "SubLayer",
	9999: "Fork",
}

var NetStrToId = map[string]uint32{
	"LivePublic": 175,
	"TestNet":    199,
	"Local":      0,
}

func NetworkIdResolver(id uint32) string {
	n, ok := Networks[id]
	if ok {
		return n
	}
	return strconv.Itoa(int(id))
}

// chain implements Chain interface
type chain struct {
	name      string
	NetworkId uint32
	Version   string
	genesis   GenesisBlock
	params    map[string]any
	Bootnodes []string
}

func (c *chain) Name() string {
	return c.name
}

func (c *chain) Genesis() GenesisBlock {
	return c.genesis
}

func (c *chain) Params() map[string]any {
	return c.params
}

func (c *chain) BootNodes() []string {
	return c.Bootnodes
}

func (c *chain) InitializeNewBlockchain(networkName string) Chain {

	ch := &chain{
		name:      "Quantos",
		NetworkId: NetStrToId[networkName],
		Version:   Version.String(),
		genesis:   c.InitializeGenesis(),
	}
	ch.params["chainId"] = NetStrToId
	ch.params["vm"] = "qvm"
	ch.params["consensus"] = []string{"qdpos", "pod"}

	return ch

}

func (c *chain) InitializeGenesis() GenesisBlock {
	genesis := &blockchain.GenesisBlock{}

	return genesis
}

type ChainImporter struct {
	ChainGetter
}
