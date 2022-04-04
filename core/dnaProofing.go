package core

import (
	"github.com/quantosnetwork/dev-0.1.0/crypto/pod"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
)

type DNAProofing interface {
	GetProof(block *pb.Block) (string, string)
	Reset() *pod.TemperProofParams
}

type POD struct {
	proofs map[*pb.BlockHeader]string
	hmacs  map[*pb.BlockHeader]string
}

func ReBuildBlockchainProofs(items *pb.Blockchain) DNAProofing {
	P := &POD{
		proofs: map[*pb.BlockHeader]string{},
		hmacs:  map[*pb.BlockHeader]string{},
	}
	for _, item := range items.GetBlocks() {
		header := item.GetHead()
		proof, hmac := P.GetProof(item)
		P.proofs[header] = proof
		P.hmacs[header] = hmac
	}

	return P

}

func (p *POD) GetProof(block *pb.Block) (string, string) {
	tp := pod.TemperProofParams{
		MutationRate:   0.008,
		ParseDuration:  15,
		PopulationSize: 5000,
		MaxFitness:     1000,
	}
	return pod.GetProof(tp, block.GetHead().GetTimestamp().AsTime())
}

func (p POD) Reset() *pod.TemperProofParams {
	return &pod.TemperProofParams{
		MutationRate:   0,
		ParseDuration:  0,
		PopulationSize: 0,
		MaxFitness:     0,
	}
}
