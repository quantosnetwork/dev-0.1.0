package module

import (
	"errors"
	"fmt"
	"github.com/looplab/fsm"
	"github.com/quantosnetwork/dev-0.1.0/color"
)

type ErroredModule struct {
	name  string
	error error
	state string
	event string
}

func (em ErroredModule) Info() string {
	inf := fmt.Sprintf("name: %s\n state: %s\n eventname: %s\n error message: %s", em.name, em.state, em.event,
		em.error)
	return inf
}

func (em ErroredModule) String() string {
	return em.Info()
}

var LoadingErrors map[string]error = map[string]error{
	"FAILED_VERIFICATION": errors.New(fmt.Sprintf(color.Red+"(INVALID MODULE) Info: %s, failed verification",
		ErroredModule.Info)),
	"UNDEFINED_MODULE_ERROR": errors.New(fmt.Sprintf(color.Red+"(UNKNOWN ERROR) : %s", ErroredModule.Info)),
}

// ModuleRepository is holding all the modules for the application, they are unsorted
// Modules map elements as map => module.Name : *ModuleRepositoryItem
type ModuleRepository struct {
	Modules map[string]*ModuleRepositoryItem
}

type ModuleRepositoryItem struct {
	Module
	ModuleFunc func(...any) any
	LoadOrder  int
	Sorted     bool
	Loaded     bool
	Errored    bool
	Stats      any
}

type LoadedModules struct {
	// State implements the final state manager provided by the looplab fsm golang 3rd party library
	State *fsm.FSM
}

// mmr is the module repository manager interface
type mmr interface {
	VerifyAllOfType(modType string) error
	ModExists(modType string, modName string) bool
	LoadBefore() (any, error)
	LoadOnlyType(modType string) (any, error)
	// LoadOnce will load only one instance of the module (singleton pattern)
	LoadOnce() (any, error)
}

type moduleInitializer interface {
	initState() *fsm.FSM
	isInitialized() bool
	preload() error
}

type modulePreloader interface {
	// SortedModules will get all modules sorted in the right load order based on priority
	SortedModules() OrderedModulesByPriority
	// WalkAndVerifyIntegrity will return the list of modules true if ok | false if integrity test fails
	WalkAndVerifyIntegrity(list OrderedModulesByPriority) map[string]bool
}

type OrderedModulesByPriority []Module

type SortByPriority []Module

func (sbp SortByPriority) Len() int               { return len(sbp) }
func (sbp SortByPriority) Swap(i int, j int)      { sbp[i], sbp[j] = sbp[j], sbp[i] }
func (sbp SortByPriority) Less(i int, j int) bool { return sbp[i].Priority < sbp[j].Priority }

func Loader(mmr) {

	//sort repo

}
