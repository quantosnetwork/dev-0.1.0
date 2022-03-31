package trie

import (
	"encoding/json"
	"github.com/quantosnetwork/dev-0.1.0/fs"
	"github.com/wealdtech/go-merkletree"
	"strconv"
	"time"
)

type QMerkleTree interface {
	NewMerkleTree(data [][]byte)
	NewMerkleFromStruct(data []any)
}

type Proof struct {
	Index     int
	Proof     string
	Height    int
	Timestamp int64
}

type MTree struct {
	ChainID      string
	ChainVersion string
	Tree         *merkletree.MerkleTree
	Proofs       map[int]*Proof
}

func (mt MTree) NewMerkleTree(data [][]byte) {
	tree, err := merkletree.New(data)
	if err != nil {
		panic(err)
	}
	mt.Tree = tree

}

func (mt MTree) NewMerkleFromStruct(data []any) {

}

func WriteTree(mt MTree) {
	F := fs.NewFileSystem()
	bt, _ := json.Marshal(mt)
	F.WriteFile(mt.GetFilePath(), bt, 0600)
}

func (mt MTree) GetFilePath() string {
	strTime := strconv.Itoa(int(time.Now().UnixNano()))
	path := "./data/" + mt.ChainID + "/" + mt.ChainVersion + "_" + strTime + "/snapshots/tree.trie"
	return path
}

func NewMerkleTree(chainID, chainVersion string, data ...[]byte) *MTree {

	bufData := make([][]byte, len(data)+1)
	for i, d := range data {
		bufData[i] = d
	}
	mt := &MTree{}
	mt.ChainID = chainID
	mt.ChainVersion = chainVersion
	mt.NewMerkleTree(bufData)

	Proofs := map[int]*Proof{}
	mt.Proofs = Proofs
	//WriteTree(mt)
	return mt

}
