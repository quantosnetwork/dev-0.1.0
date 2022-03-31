package blockchain

import (
	"encoding/json"
	"github.com/quantosnetwork/dev-0.1.0/factory"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
)

type BlockFactory interface {
	CreateEmpty() interface{}
	CreateGenesis() *GenesisBlock
	Set(key string, value interface{})
	Validate() bool
	Build() *pb.Block
}

type blockFactory struct {
	block *pb.Block
	bmap  map[string]interface{}
}

func (b blockFactory) CreateGenesis() *GenesisBlock {
	//TODO implement me
	return b.CreateEmpty().(*GenesisBlock)
}

func (b blockFactory) CreateEmpty() interface{} {
	return factory.GetFactory().CreateEmpty(b)
}

func (b blockFactory) Set(key string, value interface{}) {
	b.bmap = factory.GetFactory().ConvertTypeToMap(b.block)
	b.bmap[key] = value
}

func (b blockFactory) Validate() bool {
	//TODO implement me
	return true
}

func (b blockFactory) Build() *pb.Block {
	mapToBytes, _ := json.Marshal(b.bmap)
	bb, _ := factory.GetFactory().BuildFromBytes(mapToBytes, b.block)
	return bb.(*pb.Block)
}

func GetBlockFactory(bf BlockFactory) BlockFactory {
	eb := bf.CreateEmpty()
	return blockFactory{block: eb.(*pb.Block)}

}
