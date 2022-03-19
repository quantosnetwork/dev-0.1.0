package nodes

import (
	"github.com/quantosnetwork/v0.1.0-dev/config"
)

type Master struct {
	nodeConfig config.NodeConfig
	chain      config.ChainConfig
}
