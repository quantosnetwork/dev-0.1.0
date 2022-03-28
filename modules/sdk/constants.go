package sdk

import (
	"github.com/quantosnetwork/dev-0.1.0/version"
	"time"
)

type NetworkID byte

const (
	LIVE NetworkID = iota + 0x0ba
	TEST
	LOCAL
)

var (
	COIN_NAME                    = "Qbit"
	COIN_SYMBOL                  = "QBX"
	COIN_MAX_AVAILABLE           = 1000000000
	COIN_GENESIS_REWARD          = 50000000
	COIN_PER_BLOCK               = 1000
	DEFLATION_RATE_START         = -9.00
	DECREASE_DEFLATION_PER_YEAR  = 0.01
	BURN_RATE                    = 9.00
	BURN_RATE_INC_YEARLY         = 0.10
	FEE_RATE_VALUE_MUL           = 100
	FEE_RATE_DIV1                = 1024
	FEE_RATE_DIV2                = 2
	BLOCK_CREATED_EVERY, _       = time.ParseDuration("10 minutes")
	BLOCK_MAX_VALIDATION_TIME, _ = time.ParseDuration("5 minutes")
)

type Version struct{ version.SemVer }
