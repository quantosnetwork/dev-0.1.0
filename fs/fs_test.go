package fs

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

type MockData struct {
	Data string
}

func loadData() []byte {
	data, _ := ioutil.ReadFile("./mockdata.json")
	return data
}

func TestFileSystem(t *testing.T) {
	t.Run("test_filesystem", func(t *testing.T) {
		TestFilesystem_LoadToMemory(t)
	})
}

func TestFilesystem_LoadToMemory(t *testing.T) {
	b := NewFS(osfs.New("."))
	a := assert.New(t)
	_, err := b.LoadToMemory("test")
	if err != nil {
		t.Error(err)
	}

	f, err := b.io.OpenFile("./test/mockdata.json", os.O_RDONLY, 0600)
	if err != nil {
		t.Error(err)
	}
	s, _ := b.io.Stat("./test/mockdata.json")
	buf := make([]byte, s.Size())
	_, err = f.Read(buf)
	if err != nil {
		t.Error(err)
	}
	spew.Dump(buf)

	a.NotEmpty(buf)

}

type d struct {
	Data string
}

func TestFilesystem_SyncToDisk(t *testing.T) {

	a := NewFS(osfs.New("."))
	dd := &d{"data"}
	data, _ := json.Marshal(dd)
	util.WriteFile(a.io, "./test/mockdata.json", data, 0600)

}

func TestNewWriter(t *testing.T) {
	mock := &MockData{
		Data: "test",
	}
	w := NewWriter(WriteWithFromStruct(mock))
	err := w.Write("test.test", "testing data test writer")
	if err != nil {
		t.Fatal(err)
	}
}
