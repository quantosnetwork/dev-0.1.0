package fs

import "github.com/spf13/afero"

type QFile struct {
	f           afero.File
	flags       []uint32
	txId        string
	signature   string
	lockState   bool
	blockHeight uint32
}

type FileFlags map[uint32]string

var flags FileFlags = map[uint32]string{
	0: "nil",
	1: "default",
	2: "encrypted",
	3: "locked",
	4: "unlocked",
	5: "signed",
	6: "hasChecksum",
	7: "noUser",
	8: "immutable",
}

const (
	FLAG_NIL uint32 = iota
	FLAG_DEFAULT
	FLAG_ENC
	FLAG_LOCK
	FLAG_UNLOCK
	FLAG_SIGNED
	FLAG_CHKSUM
	FLAG_NOACCESS
	FLAG_IMMUTABLE
)

type FileManager struct{ work *QFile }

func (f FileManager) New(file afero.File, hasflags bool, params map[string]interface{}) *QFile {
	f.work = new(QFile)
	f.work.f = file
	if hasflags {
		for _, flag := range params["flags"].([]uint32) {
			f.work.flags = append(f.work.flags, flag)
		}
	} else {
		f.work.flags = append(f.work.flags, 0)
	}
	f.work.txId = params["tx_id"].(string)
	f.work.blockHeight = params["block_height"].(uint32)
	f.work.signature = params["signature"].(string)
	f.work.lockState = params["lock"].(bool)
	return f.work
}

func (f *qFs) WriteFileWithFlags(file *QFile) error {
	return nil
}
