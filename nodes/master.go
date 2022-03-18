package nodes

import (
	"quantos/0.1.0/config"
)

type Master struct {
	nodeConfig config.NodeConfig
	chain      config.ChainConfig
}
