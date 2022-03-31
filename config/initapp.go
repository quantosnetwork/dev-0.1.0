package config

import (
	"github.com/golang/protobuf/proto"
	"github.com/quantosnetwork/dev-0.1.0/cmd"
	chain "github.com/quantosnetwork/dev-0.1.0/core/blockchain"
	"io/ioutil"
	"sync"
)

func Initializer() error {

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)
	// InitializePaths creates missing paths

	go InitializePaths()
	go InitDB()
	go func() {
		bc := chain.NewBlockchain(chain.LIVE_NETWORK)
		bcj, _ := proto.Marshal(bc)
		ioutil.WriteFile("./chain.json", bcj, 0600)
		return
	}()

	wg.Wait()
	cmd.Execute()

	return nil
}
