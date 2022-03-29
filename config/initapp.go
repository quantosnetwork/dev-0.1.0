package config

import (
	"github.com/quantosnetwork/dev-0.1.0/cmd"
	"sync"
)

func Initializer() error {

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)
	// InitializePaths creates missing paths

	go InitializePaths()
	go InitDB()

	cmd.Execute()

	return nil
}
