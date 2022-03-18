package tx

type TxManager interface {
	CreateTransaction(from, to, amount string, currencyID uint)
	ValidateTransaction(validator, txid string) bool
	IncludeTransactionInBlock(txid string, txdata string) bool
}
