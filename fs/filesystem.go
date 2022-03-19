package fs

import (
	"github.com/spf13/afero"
)

type iFs interface {
	InitMemFs()
	InitOSFs()
	GetMemFs() *afero.Fs
	GetOsFs() *afero.Fs
	GetQfs() *qFs
}

type qFs struct {
	afero.Afero
	mem afero.Fs
	os  afero.Fs
}

func NewFileSystem() *qFs {
	QFS := new(qFs)
	QFS.InitMemFs()
	QFS.InitOsFs()

	return QFS
}

func (f *qFs) InitMemFs() {
	f.mem = afero.NewMemMapFs()
}

func (f *qFs) InitOsFs() {
	f.os = afero.NewOsFs()
}

func (f *qFs) GetMemFs() afero.Fs {
	return f.mem
}

func (f *qFs) GetOsFs() afero.Fs {
	return f.os
}

func (f *qFs) GetQfs() *qFs {
	return f
}
