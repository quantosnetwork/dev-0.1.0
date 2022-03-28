package tx

type TxManager interface {
	CreateTransaction(from, to, amount string, currencyID uint) (*Transaction, error)
	ValidateTransaction(validator, txid string) bool
	IncludeTransactionInBlock(txid string, txdata string) bool
}
