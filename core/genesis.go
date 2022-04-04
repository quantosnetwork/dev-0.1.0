package core

import (
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GenesisBlock interface {
	Config() *any
	Nonce() [8]byte
	Timestamp() timestamppb.Timestamp
	Pb() *pb.Block
	Raw() []byte
	DNAProof() string
	Hash() any
	Hex() string
	String() string
	Coinbase() string
	Metadata() map[string]any
}

type genesis struct{}

func (g genesis) Config() *any {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Nonce() [8]byte {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Timestamp() timestamppb.Timestamp {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Pb() *pb.Block {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Raw() []byte {
	//TODO implement me
	panic("implement me")
}

func (g genesis) DNAProof() string {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Hash() any {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Hex() string {
	//TODO implement me
	panic("implement me")
}

func (g genesis) String() string {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Coinbase() string {
	//TODO implement me
	panic("implement me")
}

func (g genesis) Metadata() map[string]any {
	//TODO implement me
	panic("implement me")
}

func CreateNewGenesisBlock() GenesisBlock {
	return &genesis{}
}
