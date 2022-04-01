package protocol

import (
	"github.com/hashicorp/go-hclog"
)

type ProtocolServer struct {
	logger hclog.Logger
	config *Config
}
