package version

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/sha3"
)

type SemVer [3]int

func (v SemVer) String() string {
	verString := fmt.Sprintf("v%x.%x.%x", v[0], v[1], v[2])
	return verString
}

func (v *SemVer) Set(major, minor, patch int) {
	v[0] = major
	v[1] = minor
	v[0] = patch
}

func (v SemVer) Hash() []byte {
	hashes := sha3.New256()
	hashes.Write([]byte(v.String()))
	return hashes.Sum(nil)
}

func (v SemVer) Verify(other Version) bool {
	return bytes.Compare(v.Hash(), other.Hash()) == 0
}

func (v SemVer) Get() SemVer {
	return v
}

type Version interface {
	String() string
	Set(major, minor, patch int)
	Hash() []byte
	Verify(other Version) bool
	Get() SemVer
}
