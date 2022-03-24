package sdk

import (
	"encoding/json"
	"io/ioutil"
)

type GenesisBlock struct {
}

/*
GenesisData

 "Coin": {
    "name": "QBit",
    "symbol": "QBX",
    "maxAvailable": "1000000000",
    "genesisReward": "50000000",
    "blockReward": "1000",
    "minUnit": "0.000000000000000001",
    "maxUnit": "10",
    "minUnitName": "QBitly",
    "maxUnitName": "DQBit",
    "coinbaseAddress": "000000004c0ac56e015b9bd5912508e1c4a992"
  },
  "Validators": [
    {
      "address": "",
      "stake": "",
      "account": ""
    }
  ],
  "Bootstrap": [{
    "host": "",
    "port": ""
  }]
*/

type GenesisData struct {
	chainID             string `json:"ChainID"`
	network             int    `json:"network"`
	payload             *GenesisPayload
	stableGasFeePercent float32 `json:"StableGasFeePercent"`
	validators          []*GenesisValidator
	coin                *Coin
}

type GenesisPayload struct {
	genesisBytes []byte `json:"GenesisByte"`
}

type GenesisValidator struct {
	Address       string
	StakingAmount float64
}

type BootstrapNode struct {
	host string
	port string
}

func NewFromGenesisData() (*Coin, *GenesisData, error) {
	loaded, err := ioutil.ReadFile("./genesisdata.json")
	if err != nil {
		return nil, nil, err
	}
	var g GenesisData
	err = json.Unmarshal(loaded, &g)
	if err != nil {
		return nil, nil, err
	}
	coin := &Coin{}
	return coin, &g, nil
}
