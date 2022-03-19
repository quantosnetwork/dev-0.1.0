package blockchain

import (
	"encoding/json"
	"github.com/quantosnetwork/v0.1.0-dev/factory"
)

type BlockFactory interface {
	CreateEmpty() interface{}
	CreateGenesis() *BlockV1
	Set(key string, value interface{})
	Validate() bool
	Build() *Block
}

type blockFactory struct {
	block *BlockV1
	bmap  map[string]interface{}
}

func (b blockFactory) CreateGenesis() *BlockV1 {
	//TODO implement me
	return b.CreateEmpty().(*BlockV1)
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

func (b blockFactory) Build() *Block {
	mapToBytes, _ := json.Marshal(b.bmap)
	bb, _ := factory.GetFactory().BuildFromBytes(mapToBytes, b.block)
	return bb.(*Block)
}

func GetBlockFactory(bf BlockFactory) BlockFactory {
	eb := bf.CreateEmpty()
	return blockFactory{block: eb.(*BlockV1)}

}
