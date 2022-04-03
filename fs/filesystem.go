package fs

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/google/uuid"
	"io"
	"os"
)

type Filesystem struct {
	io billy.Filesystem
	billy.Basic
	billy.Dir
	billy.TempFile
	billy.Chroot
	billy.Change
}

var fs billy.Filesystem

func NewFS(fs billy.Filesystem) Filesystem {
	s := Filesystem{io: fs}
	return s
}

func GetFileSystem(basePath string) Filesystem {
	return NewFS(osfs.New(basePath))
}

func (f Filesystem) LoadToMemory(path string) (*memfs.Memory, error) {
	mem := new(memfs.Memory)
	f.io = osfs.New(".")
	err := f.io.MkdirAll(path, 0755)
	if err != nil {
		return nil, nil
	}
	files, err := f.io.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src, err := f.io.Open(file.Name())
		if err != nil {
			return nil, err
		}

		dst, err := f.io.Create(file.Name())
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(dst, src); err != nil {
			return nil, err
		}

		if err := dst.Close(); err != nil {
			return nil, err
		}

		if err := src.Close(); err != nil {
			return nil, err
		}
	}
	return mem, nil
}

func (f Filesystem) SyncToDisk(onDiskPath string, data *memfs.Memory) error {
	f.io = osfs.New(".")

	dst, err := f.io.Create(onDiskPath)
	if err != nil {
		return err
	}
	src, err := data.Open(onDiskPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	if err := dst.Close(); err != nil {
		return err
	}

	if err := src.Close(); err != nil {
		return err
	}

	return nil

}

func (f *Filesystem) TmpToDisk(onDiskPath string, tmpPath string) error {
	f.io = osfs.New(".")
	exists := f.FileExists(onDiskPath)
	if !exists {
		osf := osfs.New(f.io.Join(f.io.Root(), "/"))
		dst, err := osf.Create(onDiskPath)
		if err != nil {
			return err
		}
		src, err := f.io.TempFile(tmpPath, uuid.New().String())
		if err != nil {
			return err
		}
		buf := make([]byte, 0, 4096)
		var err2 error
		_, err2 = src.Read(buf)
		if err2 != nil {
			return err
		}
		_, err2 = dst.Write(buf[:])
		if err2 != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}
		if err := dst.Close(); err != nil {
			return err
		}

		if err := src.Close(); err != nil {
			return err
		}

	}
	return nil
}

func (f Filesystem) WriteToTmp(tmpPath string, data []byte) error {
	tmp, _ := f.io.TempFile(tmpPath, uuid.New().String())
	_, err := tmp.Write(data)
	return err
}

func (f Filesystem) FileExists(onDiskPath string) bool {
	f.io = osfs.New(".")
	_, exists := f.io.OpenFile(onDiskPath, os.O_RDONLY, 0600)
	if exists != nil {
		return true
	}
	return false
}
