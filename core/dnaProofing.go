package core

import (
	"github.com/quantosnetwork/dev-0.1.0/crypto/pod"
)

type DNAProofing interface {
	GetOperator() *pod.DnaOperator
}
