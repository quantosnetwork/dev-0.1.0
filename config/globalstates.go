package config

import (
	"fmt"
	"github.com/looplab/fsm"
	"github.com/quantosnetwork/dev-0.1.0/core/blockchain"
	"github.com/quantosnetwork/dev-0.1.0/logger"
	"log"
)

// Global States configuration

type GlobalState uint32

const (
	STATE_IDLE GlobalState = iota
	STATE_INITIALIZING
	STATE_INITIALIZED
	STATE_ERROR
	STATE_READY
	STATE_CREATING_GENESIS
	STATE_GENESIS_DONE
)

var GS map[GlobalState]string = map[GlobalState]string{
	STATE_IDLE:             "idle",
	STATE_INITIALIZING:     "initializing",
	STATE_INITIALIZED:      "initialized",
	STATE_ERROR:            "error",
	STATE_READY:            "ready",
	STATE_CREATING_GENESIS: "create_genesis",
	STATE_GENESIS_DONE:     "genesis_done",
}

type GlobalStateMachine struct {
	*fsm.FSM
}

func NewGlobalStateMachine() *GlobalStateMachine {
	newFSM := fsm.NewFSM(
		GS[STATE_INITIALIZING],
		fsm.Events{
			{Name: "initializing", Src: []string{GS[STATE_INITIALIZING]}, Dst: GS[STATE_INITIALIZED]},
			{Name: "initialized", Src: []string{GS[STATE_INITIALIZING]}, Dst: GS[STATE_INITIALIZED]},
			{Name: "create_genesis", Src: []string{GS[STATE_INITIALIZED]}, Dst: GS[STATE_CREATING_GENESIS]},
			{Name: "genesis_done", Src: []string{GS[STATE_CREATING_GENESIS]}, Dst: GS[STATE_GENESIS_DONE]},
			{Name: "ready", Src: []string{GS[STATE_GENESIS_DONE]}, Dst: GS[STATE_READY]},
		},
		fsm.Callbacks{
			GS[STATE_INITIALIZING]: func(e *fsm.Event) {
				err := Initializer()
				if err == nil {
					e.FSM.SetMetadata("error", false)
					log.Println("Quantos is initializing....")
				}

			},
			GS[STATE_INITIALIZED]: func(e *fsm.Event) {
				_, ok := e.FSM.Metadata("error")
				if !ok {
					e.FSM.SetMetadata("initialized", true)
					e.FSM.SetMetadata("creating_genesis", true)
					log.Println("Quantos is successfully initialized")
				}
			},
			GS[STATE_CREATING_GENESIS]: func(e *fsm.Event) {
				log.Println("creating genesis block....")
				g := blockchain.NewLiveGenesisBlock()
				err := blockchain.WriteGenesisBlock(g.NetworkID, g)
				if err != nil {
					return
				}
				e.FSM.SetMetadata("genesis", g)

			},
			GS[STATE_GENESIS_DONE]: func(e *fsm.Event) {
				log.Println("genesis block successfully created")
			},
		})
	return &GlobalStateMachine{newFSM}
}

var InitState *GlobalStateMachine

func init() {
	InitState = NewGlobalStateMachine()
	s := GetCurrentGlobalState()
	log.Println(s)
	err := InitState.Event("initializing")

	if err != nil {
		fmt.Println(err.Error())
	} else {

		PrintState(GetCurrentGlobalState())
	}
}

func GetCurrentGlobalState() string {
	return InitState.Current()
}

func PrintState(s string) {

	logger.Logger.Info("state " + s)

}
