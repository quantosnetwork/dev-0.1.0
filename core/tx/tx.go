package tx

import (
	"github.com/mr-tron/base58"
	"github.com/quantosnetwork/v0.1.0-dev/config"
	"go.uber.org/atomic"
	"lukechampine.com/frand"
	"strconv"
	"sync"
)

type Transaction struct {
	ID         string
	MerkleRoot string
	From       string
	TxType     Type
	Time       *config.Time
	Timestamp  int64
	Fees       float64
	Nonce      uint32
	Validators []string
	Inputs     []map[string]*Input
	Outputs    []map[string]*Output
}

type Input struct {
	ID        string
	From      string
	Recipient string
	Amount    float64
	Outputs   []*Output
}

type Output struct {
	ID            string
	Address       string
	Spent         atomic.Bool
	Confirmations atomic.Int32
	Signatures    []map[string]string
}

type TXQueue struct {
	items sync.Map
}

func (t *Transaction) CreateTransaction(from, to, amount string, currencyID uint) (*Transaction, error) {
	t.From = from
	t.Inputs = []map[string]*Input{}
	inputID := generateUniqID()
	f, _ := strconv.ParseFloat(amount, 64)
	t.Inputs[0][inputID] = &Input{
		ID:        inputID,
		From:      from,
		Recipient: to,
		Amount:    f,
	}
	t.Outputs = []map[string]*Output{}
	outputID := generateUniqID()
	t.Outputs[0][outputID] = &Output{
		ID:            outputID,
		Address:       to,
		Spent:         atomic.Bool{},
		Confirmations: atomic.Int32{},
		Signatures:    []map[string]string{},
	}
	t.Outputs[0][outputID].Spent.Store(false)
	t.Outputs[0][outputID].Confirmations.Store(0)
	return t, nil

}

func generateUniqID() string {
	buf := make([]byte, 32)
	frand.Read(buf)
	return base58.Encode(buf)
}
