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
