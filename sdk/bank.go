package sdk

import (
	"github.com/google/uuid"
	"go.uber.org/atomic"
)

type Coin struct {
	ID              string
	Name            string  `json:"coin.name"`
	Symbol          string  `json:"coin.symbol"`
	MaxAvailable    float64 `json:"coin.maxAvailable"`
	GenesisReward   float64 `json:"coin.genesisReward"`
	BlockReward     float64 `json:"coin.blockReward"`
	Decimals        int     `json:"coin.decimals"`
	CoinbaseAddress string  `json:"coin.coinbaseAddress"`
}

type Token struct {
	ID                string
	Name              string
	Symbol            string
	Decimals          int
	MaxAvailable      float64
	Mintable          bool
	Burnable          bool
	Upgradeable       bool
	Killable          bool
	Tradable          bool
	GasFee            float64
	CreationFee       float64
	CreationTxAddress string
	BaseCoin          *Coin
	CreatorAddress    string
	ContractAddress   string
	ContractABI       string
	OwnerAddress      string
	ContractCode      string
	TokenType         string
	TokenStruct       any
	ApprovalSignature string
	ApprovedBy        string
	Notes             []string
}

type CoinUnit struct {
	ID          string
	MintedOn    int64
	MintedBy    string
	MintedFor   string
	MintedTx    string
	DnaProof    string
	Owner       string
	Value       float64
	Valid       bool
	Blacklisted bool
	Locked      bool
	Spendable   bool
}

type Bank struct {
	Address        string
	Accounts       map[string]*accountManager
	ElectedComitee []*accountManager
	BaseCurrency   *Coin
	OtherAssets    []*Coin
	Liquidities    map[string]interface{}
	Circulating    map[string]*CoinUnit
	DeadCoins      map[string]*CoinUnit
	NumCirculating atomic.Uint64
	TotalAvailable atomic.Uint64
	PreMinted      uint64
	Minted         uint64
	LeftToMint     uint64
	minter         BaseMinter
}

func (b *Bank) GetBaseBank() *Bank {
	return b
}

func (b *Bank) GetBankAddress() string {
	return b.Address
}

func (b *Bank) GetBankAccount(id string) *accountManager {
	return b.Accounts[id]
}

func (b *Bank) GetElectedCommitee() []*accountManager {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) TotalCoinsAvailable() uint64 {
	return b.TotalAvailable.Load()
}

func (b *Bank) TotalCirculating() uint64 {
	return b.Minted
}

func (b *Bank) GetMinter() BaseMinter {
	return b.minter
}

func (b *Bank) IssueNewAsset(c *Coin) error {
	b.OtherAssets = append(b.OtherAssets, c)
	return nil
}

func (b *Bank) BuyLiquidities(amount uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) SellLiquidities(amount uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) TransferFrom(fromAcct string, toAcct string, amount float64) {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) SignTransfer(txID string) {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) CancelTransfer(txID string) {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) SetTransferState(txID string, stateID int) {
	//TODO implement me
	panic("implement me")
}

func (b *Bank) CreateNewBankAccount() *accountManager {
	//TODO implement me
	panic("implement me")
}

type BaseMinter interface {
	mint(currencyID string, amount uint64)
	burn(currencyID string, amount uint64)
	info() map[string]interface{}
}

type BankManager interface {
	NewBaseCoin(info map[string]any)
	GetBaseBank() *Bank
	GetBankAddress() string
	GetBankAccount(id string) *accountManager
	GetElectedCommitee() []*accountManager
	TotalCoinsAvailable() uint64
	TotalCirculating() uint64
	GetMinter() BaseMinter
	IssueNewAsset(c *Coin) error
	BuyLiquidities(amount uint64) error
	SellLiquidities(amount uint64) error
	TransferFrom(fromAcct string, toAcct string, amount float64)
	SignTransfer(txID string)
	CancelTransfer(txID string)
	SetTransferState(txID string, stateID int)
	CreateNewBankAccount() *accountManager
}

func (b *Bank) NewBaseCoin(info map[string]any) {

	c := &Coin{}
	c.ID = uuid.New().String()
	c.Decimals = info["decimals"].(int)
	c.Name = info["name"].(string)
	c.Symbol = info["symbol"].(string)
	c.MaxAvailable = info["maxAvailable"].(float64)
	c.GenesisReward = info["genesisReward"].(float64)
	c.BlockReward = info["blockReward"].(float64)
	c.CoinbaseAddress = info["coinbaseAddress"].(string)
	b.BaseCurrency = c
}

func NewGlobalBank() BankManager {

	bank := new(Bank)
	bank.minter = NewMinter()

	return bank

}

func NewMinter() BaseMinter {
	return &minter{}
}

type minter struct {
	BaseMinter
}

func (m *minter) mint(currencyID string, amount uint64) {

}

func (m *minter) burn(currencyID string, amount uint64) {

}

func (m *minter) info() map[string]interface{} {
	return map[string]interface{}{}
}
