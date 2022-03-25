package p2p

type Keys interface {
	GenerateNewKeyPair()
	LoadFromDisk()
	WriteToDisk()
}

type SignKey interface {
	Keys
	Sign(msg interface{})
	VerifySignature()
}

type EncryptionKey interface {
	Keys
	Encrypt(msg interface{})
	Decrypt(msg interface{})
}

type Key interface {
	String() string
	Bytes() ([]byte, error)
	Hex() string
}

type PublicKey Key

type PrivateKey Key
