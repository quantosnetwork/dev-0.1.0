package sdk

import (
	"github.com/quantosnetwork/v0.1.0-dev/version"
)

type NetworkID uint32

const (
	LIVE NetworkID = iota + 1
	TEST
	LOCAL
)

type Version struct{ version.SemVer }
