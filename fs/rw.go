package fs

import (
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-billy/v5/util"
)

type ReadRequest struct {
	encoding  string
	encrypted bool
	key       string
	decrypted []byte
	toStruct  bool
	Struc     *any
}

type WriteRequest struct {
	encoding   string
	encrypted  bool
	fromStruct bool
	Struc      *any
	key        string
}

type ReadOption func(r *Reader)
type WriteOption func(w *Writer)

func (rq *ReadRequest) DefaultReadOptions() *ReadRequest {
	return &ReadRequest{
		encrypted: false,
		encoding:  "json",
		key:       "",
		toStruct:  true,
	}
}

type Reader struct {
	fs       Filesystem
	Filename string
	Options  ReadRequest
	buffer   []byte
	IReader
}

type Writer struct {
	fs       Filesystem
	Filename string
	Options  WriteRequest
	IWriter
}

func (read *Reader) Read() error {

	file, err := read.fs.io.Open(read.Filename)
	if err != nil {
		return err
	}
	s, _ := read.fs.io.Stat(read.Filename)
	read.buffer = make([]byte, s.Size())
	_, err = file.Read(read.buffer)
	if err != nil {
		return err
	}
	return nil
}

func NewReader(opts ...ReadOption) *Reader {
	reader := &Reader{}
	for _, opt := range opts {
		opt(reader)
	}
	return reader
}

func ReadWithToStruct(struc any) ReadOption {
	return func(r *Reader) {
		r.Options.Struc = &struc
		r.Options.toStruct = true
	}
}

func WriteWithFromStruct(struc any) WriteOption {
	return func(w *Writer) {
		w.Options.fromStruct = true
		w.Options.Struc = &struc
	}
}

func ReadEncrypted(filename string) (string, error) {
	r := NewReader(nil)
	r.Filename = filename
	err := r.Read()
	if err != nil {
		return "", err
	}
	buf := r.buffer
	return r.Decrypt(string(buf))
}

func ReadEncryptedToStruct(filename string, s any) (any, error) {
	r := NewReader(ReadWithToStruct(s))
	r.Filename = filename
	err := r.Read()
	if err != nil {
		return nil, err
	}
	out, err := r.DecryptToStruct(string(r.buffer))
	return out, err

}

func NewWriter(opts ...WriteOption) IWriter {
	writer := new(Writer)
	writer.fs.io = osfs.New(".")
	if len(opts) > 0 && opts != nil {
		for _, opt := range opts {
			opt(writer)
		}
	}
	return writer

}

func (w *Writer) Write(filename string, data string) error {

	buffer := make([]byte, len([]byte(data)))
	copy(buffer, []byte(data))

	err := util.WriteFile(w.fs.io, filename, buffer, 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteEncrypted(filename string, data string) error {
	w := NewWriter()
	in, err := w.Encrypt(data)
	if err != nil {
		return err
	}
	err = w.Write(filename, in)
	if err != nil {
		return err
	}
	return nil
}

func NewWriterWithoutOptions() *Writer {
	writer := new(Writer)
	writer.fs.io = osfs.New(".")
	return writer
}

func NewReaderWithoutOptions() *Reader {
	reader := new(Reader)
	reader.fs.io = osfs.New(".")
	return reader
}

type FileRW interface {
	NewWriterWithoutOptions() *Writer
	NewReaderWithoutOptions() *Reader
	NewWriter(opts ...WriteOption) *Writer
	NewReader(opts ...ReadOption) *Reader
}

type IWriter interface {
	Write(filename string, data string) error
	Encrypt(data string) (string, error)
}

type IReader interface {
	Read() error
	Decrypt(data string) (string, error)
}
