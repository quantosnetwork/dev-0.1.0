package config

import (
	"github.com/quantosnetwork/dev-0.1.0/datastore"
)

var DB *datastore.Repo

func InitDB() {
	DB = datastore.NewDatastore()
}
