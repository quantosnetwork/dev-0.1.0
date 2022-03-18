package store

import (
	"encoding/json"
	"github.com/hashicorp/go-memdb"
	"io/ioutil"
	"time"
)

type MemoryDB struct {
	db     *memdb.MemDB
	schema *memdb.DBSchema
}

func NewMemoryDB(name string, tableSchemas map[string]*memdb.TableSchema) *MemoryDB {
	s := CreateMemoryDBSchema(tableSchemas)
	db, err := memdb.NewMemDB(s)
	if err != nil {
		panic(err)
	}
	return &MemoryDB{db, s}
}

func CreateMemoryDBSchema(tableSchemas map[string]*memdb.TableSchema) *memdb.DBSchema {
	schema := &memdb.DBSchema{
		Tables: tableSchemas,
	}
	return schema
}

func (m *MemoryDB) InsertOne(tableName string, data interface{}) error {
	txn := m.db.Txn(true)
	if err := txn.Insert(tableName, data); err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (m *MemoryDB) InsertMany(tableName string, data []interface{}) error {
	txn := m.db.Txn(true)
	for _, d := range data {
		if err := txn.Insert(tableName, d); err != nil {
			return err
		}
	}
	txn.Commit()
	return nil
}

func (m *MemoryDB) FindOne(tableName string, key, data string) (interface{}, error) {
	txn := m.db.Txn(false)
	defer txn.Abort()
	raw, err := txn.First(tableName, key, data)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (m *MemoryDB) SaveSnapShot() {
	snap := m.db.Snapshot()
	snapBytes, _ := json.Marshal(snap)
	ioutil.WriteFile("data/snapshots/"+time.Now().String()+".qsnap", snapBytes, 0600)
	return
}
