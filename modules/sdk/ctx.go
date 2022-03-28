package sdk

import (
	"bytes"
	"context"
	"github.com/lunixbochs/struc"
	"github.com/quantosnetwork/dev-0.1.0/config"
	"github.com/quantosnetwork/dev-0.1.0/crypto/pod"
	"github.com/quantosnetwork/dev-0.1.0/version"
	"time"
)

type QuantosBlockchainContext struct {
	Ctx            context.Context
	CurrentNetwork NetworkID
	CurrentVersion *Version
	Config         *config.ChainConfig
	GenesisDNA     string
}

func NewQuantosContext(ctx context.Context, netid NetworkID, versionMaj, versionMin,
	versionPatch int, cfg *config.ChainConfig) *QuantosBlockchainContext {

	ver := version.SemVer{}
	ver.Set(versionMaj, versionMin, versionPatch)
	qbctx := &QuantosBlockchainContext{CurrentNetwork: netid,
		CurrentVersion: &Version{ver}}
	qbctx.Config = cfg
	qbctx.Ctx = context.WithValue(ctx, []byte("qb_ctx"), qbctx)
	return qbctx
}

func (ctx *QuantosBlockchainContext) GenerateGenesisDnaProof() {
	d, _ := time.ParseDuration("5s")
	p := pod.TemperProofParams{
		MutationRate:   0.005,
		ParseDuration:  d,
		PopulationSize: 1000,
		MaxFitness:     0.5,
	}
	proof, hmac := pod.GetProof(p, time.Now())
	bbytes := make([][]byte, 2)
	bbytes[0] = []byte(proof)
	bbytes[1] = []byte(hmac)
	var buf bytes.Buffer
	struc.Pack(&buf, bbytes)
	ctx.GenesisDNA = buf.String()

}
