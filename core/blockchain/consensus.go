package blockchain

import (
	"errors"
	pb "github.com/quantosnetwork/dev-0.1.0/proto/gen/proto/quantos/pkg/v1"
	"math/big"
	"math/rand"
)

type QuantosConsensus struct {
	Round         int
	Validators    []*pb.Validator
	NumValidators int
	TotalStakes   *big.Int
	Config        *ConsensusConfig
}

type ConsensusConfig struct {
	MinValidators int
	MaxValidators int
	MinStake      *big.Int
	MinRounds     int
	MaxRounds     int
	Quorum        int
	Reward        *big.Int
}

func (c *QuantosConsensus) SetConfig() {

	minstakes := big.NewInt(5000)
	Reward := big.NewInt(5000000)

	c.Config = &ConsensusConfig{
		1, 51, minstakes, 1, 5, 30, Reward,
	}
}

func (c *QuantosConsensus) AddValidator(v *pb.Validator) {
	if c.NumValidators < c.Config.MaxValidators {
		if v.GetNode().GetStake() >= c.Config.MinStake.Uint64() {
			if c.Round <= c.Config.MinRounds {
				c.Validators = append(c.Validators, v)
			}
		}
	}
}

func (c *QuantosConsensus) GetNumValidators() int {
	return len(c.Validators)
}

func (c *QuantosConsensus) GetTotalStakes() *big.Int {

	b := new(big.Int)
	t := new(big.Int)

	for _, v := range c.Validators {
		b.SetUint64(v.GetNode().GetStake())
		t.Add(t, b)
	}
	return t
}

func (c *QuantosConsensus) NewRound() (*Round, error) {
	if c.Round == 0 {
		c.Round = 1
	}
	if c.Round > 0 {
		c.Round += 1
	}
	if c.Round <= c.Config.MaxRounds {
		// can start round
		r := &Round{
			ID:                  c.Round,
			Validators:          c.Validators,
			tmp:                 big.NewInt(0),
			totalStakes:         big.NewInt(0),
			reward:              c.Config.Reward,
			WinnerPool:          []*pb.Validator{},
			isLast:              false,
			distributeRemainder: false,
			number:              0,
		}
		return r, nil
	}

	return nil, errors.New("cannot start consensus, round limit exceeded")

}

type Round struct {
	ID                  int
	Validators          []*pb.Validator
	tmp                 *big.Int
	totalStakes         *big.Int
	reward              *big.Int
	WinnerPool          []*pb.Validator
	isLast              bool
	distributeRemainder bool
	number              uint64
}

func (c *QuantosConsensus) Start(r *Round) (*pb.Validator, error) {

	for _, v := range c.Validators {
		n := v.GetNode()
		if n.GetStake() > c.Config.MinStake.Uint64() {
			r.WinnerPool = append(r.WinnerPool, v)
			r.totalStakes.Add(r.totalStakes, big.NewInt(0).SetUint64(n.GetStake()))
		}
	}

	if r.WinnerPool == nil {
		return nil, errors.New("there are no validators in the network")
	}
	r.number = uint64(rand.Int63n(r.totalStakes.Int64()))

	for _, v := range c.Validators {
		r.tmp.Add(r.tmp, big.NewInt(0).SetUint64(v.GetNode().GetStake()))
		if r.number < r.tmp.Uint64() {
			return v, nil
		}
	}
	return nil, errors.New("no winner were picked")
}
