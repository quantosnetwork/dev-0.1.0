package module

type ModuleManager interface {
	LoadModule(moduleName string)
}

type Module struct {
	Name          string            `json:"mod_name"`
	Type          ModuleType        `json:"mod_type"`
	AuthorAddress string            `json:"mod_author,omitempty"`
	Version       *ModuleVersioning `json:"mod_version"`
	Priority      int
}

type ModuleVersioning struct {
	VerCtrlType    string
	VerCtrlRepoUrl string
	LatestMajor    string
	LatestMinor    string
	LatestPatch    string
	Timestamp      string
	Checksum       uint32
}

type ModuleType uint32

const (
	UNDEFINED ModuleType = iota
	SYSTEM
	CONSENSUS
	BANKING
	P2P
	STORAGE
	CONTENT
	TRIE
	ACCOUNT
	TX
	BLOCK_QUEUE
	TX_QUEUE
	COMMUNITY_LOCAL
	COMMUNITY_EXTERNAL
	CONTRACT_INTERNAL
	CONTRACT_EXTERNAL
)

type moduleTypeStrings map[ModuleType]string

var ModuleTypes moduleTypeStrings = map[ModuleType]string{
	UNDEFINED:          "mod_undefined",
	SYSTEM:             "mod_system",
	CONSENSUS:          "mod_consensus",
	BANKING:            "mod_banking",
	P2P:                "mod_p2p",
	STORAGE:            "mod_storage",
	CONTENT:            "mod_content",
	TRIE:               "mod_merkle_tree",
	ACCOUNT:            "mod_account",
	TX:                 "mod_tx",
	BLOCK_QUEUE:        "mod_block_queue",
	TX_QUEUE:           "mod_tx_queue",
	COMMUNITY_LOCAL:    "mod_community_loc",
	COMMUNITY_EXTERNAL: "mod_community_ext",
	CONTRACT_INTERNAL:  "mod_contract_internal",
	CONTRACT_EXTERNAL:  "mod_contract_external",
}
