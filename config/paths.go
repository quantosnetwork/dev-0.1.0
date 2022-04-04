package config

import (
	"errors"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/mr-tron/base58"
	"github.com/quantosnetwork/dev-0.1.0/fs"
	"github.com/spf13/afero"
	"lukechampine.com/blake3"
	"os"
	"path"
)

var PathRepo map[string]*Path

type Path struct {
	ID       string
	realPath string
	hidden   bool
}

type Paths interface {
	Get(id string) (*Path, error)
	Create(realPath string) (*Path, error)
	Exists(realPath string) bool
}

type InMemoryPath struct {
	Paths
	memfs *afero.MemMapFs
}

func (p Path) Exists(realPath string) bool {
	//_ = fs.NewFileSystem()

	//exists, _ := f.Exists(realPath)
	return true
}

func (p Path) Create(realPath string) (*Path, error) {
	id := generatePathID([]byte(realPath))
	path := &Path{}
	path.ID = id
	path.realPath = realPath
	path.hidden = false
	return path, nil
}

func (p Path) Get(id string) (*Path, error) {
	if PathRepo[id] != nil {
		return PathRepo[id], nil
	}
	return nil, errors.New("path does not exists")
}

func (p *InMemoryPath) Exists(realPath string) bool {
	f := fs.NewFS(osfs.New("."))
	_, err := f.Open(realPath)
	if err != nil {
		if err == os.ErrNotExist {
			return true
		}
	}

	return false
}

func generatePathID(b []byte) string {
	hasher := blake3.Sum256(b)
	return base58.Encode(hasher[:20])
}

type PathString string

const (
	BASE_PATH     PathString = "data"
	KEY_PATH      PathString = ".keys"
	DB_PATH       PathString = "db"
	WALLET_PATH   PathString = "wal"
	LIVENET_PATH  PathString = "live"
	TESTNET_PATH  PathString = "test"
	LOCALNET_PATH PathString = "local"
	TMP_PATH      PathString = ".tmp"
	CACHE_PATH    PathString = ".cache"
	LOG_PATH      PathString = "logs"
)

var pathArray []PathString = []PathString{KEY_PATH, DB_PATH, WALLET_PATH, LIVENET_PATH, TESTNET_PATH,
	LOCALNET_PATH,
	TMP_PATH, CACHE_PATH, LOG_PATH}

func InitializePaths() map[string]*Path {
	var P Path
	PathRepo = map[string]*Path{}
	for _, pth := range pathArray {
		realPath := path.Join(string(BASE_PATH), string(pth))
		if !P.Exists(realPath) {
			newPath, _ := P.Create(realPath)
			PathRepo[newPath.ID] = newPath
		}
	}
	for _, pth2 := range PathRepo {
		_ = os.MkdirAll(pth2.realPath, 0755)
	}
	return PathRepo
}

func GetRealPath(index string) string {
	return PathRepo[index].realPath
}

func GetIDPath(index string) string {
	return PathRepo[index].ID
}
