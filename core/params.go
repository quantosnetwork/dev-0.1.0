package core

type Params interface {
	Forks() *Forks
	ChainId() int
	VM() map[string]any
}

type Fork uint64

type Forks interface {
	Get(name string) *any
	At(block string) *any
	New(n uint64) *Fork
	Enabled() []uint64
	AutoEnableAt(blockNum uint64) error
}
