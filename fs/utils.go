package fs

import (
	"bytes"
	"github.com/spf13/afero"
	"os"
	"strings"
	"time"
)

// FsUtils interface provides many file system functions
type FsUtils interface {
	GetFS() *qFs
	GetBuf() *bytes.Buffer
	SplitFileName(filename string, sep string) (length int, parts []string)
	GetFilesFromDir(path string) (num int, files []afero.File, err error)
	GetFilesWithExt(ext string, inputFiles []afero.File) (num int, files []afero.File)
	GetFilesLargerThan(numBytes int64, inputFiles []afero.File) (num int, files []afero.File)
	EncryptWithKey(file afero.File, key []byte) error
	DecryptWithKey(file afero.File, key []byte) ([]byte, error)
	ConvertJsonToStruct(fileData []byte, jStruct any) any
	StoreInTmp(path string, fileData []byte, perm os.FileMode, expires time.Duration) error
	SyncTmpFileToDisk(tmpPath string, diskPath string) error
}

// FileUtils implements FsUtils interface
type FileUtils struct {
	buf *bytes.Buffer
	fs  *qFs
}

func (fu *FileUtils) GetFS() *qFs {
	return fu.fs
}

func (fu *FileUtils) GetBuf() *bytes.Buffer {
	return fu.buf
}

func GetFileUtils() FsUtils {
	qfs := NewFileSystem()
	buf := new(bytes.Buffer)
	futils := &FileUtils{buf, qfs}
	return futils
}

func (fu *FileUtils) GetFilesFromDir(path string) (num int, files []afero.File, err error) {
	fs := fu.GetFS()
	fsOs := fs.GetOsFs()
	_, err = fs.Exists(path)
	if err != nil {
		return 0, nil, err
	}

	fd, _ := fs.ReadDir(path)
	num = len(fd)
	files = make([]afero.File, num)
	for _, ff := range fd {

		if !ff.IsDir() {
			fff, _ := fsOs.Open(ff.Name())
			files = append(files, fff)
		}

	}

	return

}

func (fu *FileUtils) SplitFileName(filename string, sep string) (length int, parts []string) {
	parts = strings.Split(filename, sep)
	length = len(parts)
	return

}

func (fu *FileUtils) GetFilesWithExt(ext string, inputFiles []afero.File) (num int, files []afero.File) {

	num = 0
	files = make([]afero.File, 0, len(files))
	for _, file := range inputFiles {
		ns, _ := fu.SplitFileName(file.Name(), ext)
		if ns > 1 {
			files = append(files, file)
			num++

		}
	}
	return

}

func (fu *FileUtils) GetFilesLargerThan(size int64, inputFiles []afero.File) (num int, files []afero.File) {

	files = make([]afero.File, 0, len(inputFiles))
	num = 0

	for _, file := range inputFiles {
		s, _ := file.Stat()
		if s.Size() > size {
			files = append(files, file)
			num++
		}
	}

	return

}

func (fu *FileUtils) EncryptWithKey(file afero.File, key []byte) error {
	//TODO implement me
	panic("implement me")
}

func (fu *FileUtils) DecryptWithKey(file afero.File, key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (fu *FileUtils) ConvertJsonToStruct(fileData []byte, jStruct any) any {
	//TODO implement me
	panic("implement me")
}

func (fu *FileUtils) StoreInTmp(path string, fileData []byte, perm os.FileMode, expires time.Duration) error {
	//TODO implement me
	panic("implement me")
}

func (fu *FileUtils) SyncTmpFileToDisk(tmpPath string, diskPath string) error {
	//TODO implement me
	panic("implement me")
}
