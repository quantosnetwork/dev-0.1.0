package config

import (
	"fmt"
	"github.com/looplab/fsm"
)

// Global States configuration

type GlobalState uint32

const (
	STATE_IDLE GlobalState = iota
	STATE_INITIALIZING
	STATE_INITIALIZED
	STATE_ERROR
)

var GS map[GlobalState]string = map[GlobalState]string{
	STATE_IDLE:         "idle",
	STATE_INITIALIZING: "initializing",
	STATE_INITIALIZED:  "initialized",
	STATE_ERROR:        "error",
}

type GlobalStateMachine struct {
	*fsm.FSM
}

func NewGlobalStateMachine() *GlobalStateMachine {
	newFSM := fsm.NewFSM(
		GS[STATE_IDLE],
		fsm.Events{
			{Name: "initializing", Src: []string{GS[STATE_IDLE]}, Dst: GS[STATE_INITIALIZED]},
			{Name: "initialized", Src: []string{GS[STATE_INITIALIZING]}, Dst: GS[STATE_INITIALIZED]},
		},
		fsm.Callbacks{
			GS[STATE_INITIALIZING]: func(e *fsm.Event) {
				err := Initializer()
				if err != nil {
					e.FSM.SetMetadata("error", err.Error())
					fmt.Println("an error occurred")
				}
				if err == nil {
					e.FSM.SetMetadata("error", false)
					fmt.Println("Quantos is initializing")
				}
			},
			GS[STATE_INITIALIZED]: func(e *fsm.Event) {
				_, ok := e.FSM.Metadata("error")
				if ok {
					e.FSM.SetMetadata("initialized", true)
					fmt.Println("Quantos is successfully initialized")
				}
			},
		})
	return &GlobalStateMachine{newFSM}
}

var InitState *GlobalStateMachine

func init() {
	InitState = NewGlobalStateMachine()
}

func GetCurrentGlobalState() string {
	return InitState.Current()
}
