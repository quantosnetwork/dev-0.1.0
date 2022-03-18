package store

import (
	"context"
	"time"
)

type StorageType int

const (
	MEMORY StorageType = iota
	LEVELDB
	BADGER
	VAULT
	QUEUE
)

type StorageElementsType int

const (
	BLOCK StorageElementsType = iota
	BLOCKHEADER
	TRANSACTION
	QUEUEITEM
	PRIVWALLET
	CONTENT
	CONTRACT
	KEY
	PEER
	ACCOUNT
	EPHEMERAL_DATA
	CHAIN_UPGRADE
	RECEIPTS
	MERKLEROOT
	HASH
	MEDIA
	HTML
	SCRIPT
	EXCHANGEDATA
	BANKDATA
	JSON
	BYTES
	PROTO
)

type Store interface {
	InitDB(StorageType) Storage
	GetInstance(context.Context) Storage
}

type Query interface {
	FindOne(params ...interface{}) (Record, interface{}, error)
}

type Storage interface {
	CreateNew(ctx context.Context, params ...interface{}) (Store, error)
	GetStore(context.Context) Store
}

type RecordIndex struct {
	ID         string
	HasKeys    []string
	BucketName string
	Order      int32
}

type Record struct {
	SType     StorageType
	RType     StorageElementsType
	ID        string
	Data      []byte
	Immutable bool
	CreatedOn time.Time
	UpdatedOn time.Time
	DeletedOn time.Time
	Ref       *RecordIndex
}
