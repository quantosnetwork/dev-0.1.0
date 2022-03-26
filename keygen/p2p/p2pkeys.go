package p2p

import (
	"encoding/hex"
	"fmt"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
)

type P2PKeys interface {
	Keys
	KeyPair() p2pKeyPair
}

type P2PPrivateKey struct {
	Key
	scalar kyber.Scalar
	group  kyber.Group
	suite  edwards25519.SuiteEd25519
}

var curve = edwards25519.NewBlakeSHA256Ed25519()
var sha256 = curve.Hash()

type P2PSignature struct {
	r kyber.Point
	s kyber.Scalar
}

func (s P2PSignature) GetSig() P2PSignature {
	return s
}

type P2pSignature interface {
	GetPublicKey(message string) kyber.Point
	GetSig() P2PSignature
}

type P2PPublicKey struct {
	Key
	point kyber.Point
}

type p2pKeyPair struct {
	SK P2PPrivateKey
	PK P2PPublicKey
}

func (p2p *p2pKeyPair) KeyPair() p2pKeyPair {
	return p2pKeyPair{p2p.SK, p2p.PK}
}

func (p2p *p2pKeyPair) LoadFromDisk() {
	//TODO implement me
	panic("implement me")
}

func (p2p *p2pKeyPair) WriteToDisk() {
	//TODO implement me
	panic("implement me")
}

func (p2p *p2pKeyPair) GenerateNewKeyPair() {
	g := curve.Point().Base()
	k := curve.Scalar().Pick(curve.RandomStream())
	p := curve.Point().Mul(k, g)
	p2p.SK = P2PPrivateKey{scalar: k}
	p2p.PK = P2PPublicKey{point: p}
}

func (p2p *p2pKeyPair) PublicKey() P2PPublicKey {
	return p2p.PK
}

func (p2p *p2pKeyPair) PrivateKey() P2PPrivateKey {
	return p2p.SK
}

func Hash(s string) kyber.Scalar {
	sha256.Reset()
	sha256.Write([]byte(s))
	return curve.Scalar().SetBytes(sha256.Sum(nil))
}

func (priv P2PPrivateKey) Sign(message string) P2pSignature {
	g := curve.Point().Base()
	k := curve.Scalar().Pick(curve.RandomStream())
	r := curve.Point().Mul(k, g)
	e := Hash(message + r.String())
	s := curve.Scalar().Sub(k, curve.Scalar().Mul(e, priv.scalar))
	return P2PSignature{r: r, s: s}
}

func (s P2PSignature) GetPublicKey(message string) kyber.Point {
	// Create a generator.
	g := curve.Point().Base()

	// e = Hash(m || r)
	e := Hash(message + s.r.String())

	// y = (r - s * G) * (1 / e)
	y := curve.Point().Sub(s.r, curve.Point().Mul(s.s, g))
	y = curve.Point().Mul(curve.Scalar().Div(curve.Scalar().One(), e), y)

	return y
}

func (pub P2PPublicKey) String() string {
	return pub.point.String()
}

func (pub P2PPublicKey) Bytes() ([]byte, error) {
	return pub.point.MarshalBinary()
}

func (pub P2PPublicKey) Hex() string {
	b, _ := pub.Bytes()
	return hex.EncodeToString(b)
}

func (priv P2PPrivateKey) String() string {
	return priv.scalar.String()
}

func (priv P2PPrivateKey) Bytes() ([]byte, error) {
	return priv.scalar.MarshalBinary()
}

func (priv P2PPrivateKey) Hex() string {
	b, _ := priv.Bytes()
	return hex.EncodeToString(b)
}

// m: Message
// s: Signature
// y: Public key
func (pub P2PPublicKey) Verify(m string, S P2PSignature) bool {
	// Create a generator.
	g := curve.Point().Base()
	y := pub.point
	// e = Hash(m || r)
	e := Hash(m + S.r.String())

	// Attempt to reconstruct 's * G' with a provided signature; s * G = r - e * y
	sGv := curve.Point().Sub(S.r, curve.Point().Mul(e, y))

	// Construct the actual 's * G'
	sG := curve.Point().Mul(S.s, g)

	// Equality check; ensure signature and public key outputs to s * G.
	return sG.Equal(sGv)
}
func (S P2PSignature) String() string {
	return fmt.Sprintf("(r=%s, s=%s)", S.r, S.s)
}

func NewP2PKeys() P2PKeys {
	p := &p2pKeyPair{}
	p.GenerateNewKeyPair()
	return p
}

func (p2p p2pKeyPair) DerivePubKey(message string) (string, kyber.Point, kyber.Scalar) {
	S := p2p.PrivateKey().Sign(message)
	s := S.GetSig()
	DK := S.GetPublicKey(message)
	return DK.String(), s.r, s.s
}
