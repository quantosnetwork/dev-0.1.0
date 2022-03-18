package config

import (
	"github.com/looplab/fsm"
	"quantos/0.1.0/version"
)

type ChainConfig struct {
	ID          int64
	Version     version.Version
	VersionHash []byte
	Genesis     []byte
	FSM         *fsm.FSM
}

func NewChainConfig() {
	c := &ChainConfig{}
	// init default config values
	c.ID = 0
	v := c.Version.Get()
	c.VersionHash = v.Hash()
	// we initialize genesis as an empty []byte with a max len of 4096 bytes
	c.Genesis = make([]byte, 0, 4096)
	// finite state machine initialize with initial state of stopped
	c.FSM = fsm.NewFSM(
		"stopped",
		// fsm events will be declared in an external file same with callbacks
		fsm.Events{},
		fsm.Callbacks{})
}
