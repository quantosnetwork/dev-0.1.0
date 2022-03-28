package config

import (
	"sync"
)

func Initializer() error {

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)
	// InitializePaths creates missing paths
	go func() {
		InitializePaths()
	}()
	wg.Wait()
	return nil
}
