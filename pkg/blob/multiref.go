package blob

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"os"

	"github.com/multiformats/go-multihash"
)

func bug(s string) {
	fmt.Printf("%s\n", s)
}

func quickTest() {
	s := "This is the content I'm testing"
	mh := MHashFromString(s)
	h := mh.Hash()
	h.Write([]byte(s))
	if !mh.HashMatches(h) {
		fmt.Printf("QT: ERROR: We don't have a match\n")
	}

	b := RefFromString(s)
	if b != mh {
		fmt.Printf("QT: ERROR: Refs are not equal\n")
	}
}

func init() {
	quickTest()
}

// -----------------------------------------------------------
/// MY Stuff below here              -------------------------

func RefFromBytes(b []byte) Ref {
	s1 := NewHash()
	s1.Write(b)
	return RefFromHash(s1)
}

type Mhash [sha256.Size]byte

func (m Mhash) digestName() string {
	return "mh"
}

func (m Mhash) bytes() []byte {
	return m[:]
}

// ASLAK: Must implement somehow :(

func (m Mhash) newHash() hash.Hash {
	//return mhash.New()
	return sha256.New()
}

func (r Mhash) digestToString() string {
	// TODO: Ignoring the error here
	m, _ := multihash.Encode(r[:], multihash.SHA2_256)
	m2 := multihash.Multihash(m)
	return m2.B58String()
}

var sha256Meta = &digestMeta{
	ctor:  mhashFromBinary,
	ctors: mhashFromHexString,
	ctorb: mhashFromHexBytes,
	size:  sha256.Size,
}

func mhashFromBinary(b []byte) digestType {
	var d Mhash
	if len(d) != len(b) {
		panic("bogus sha-256 length")
	}
	copy(d[:], b)
	return d
}

func mhashFromHexString(hex string) (digestType, bool) {
	if len(hex) == 64 {
		fmt.Fprintf(os.Stderr, "mhashFromHexString %s length: %s\n", hex, len(hex))
		panic("64 length hex string - only happens if someone screwed up calling us")
	}
	m, er := multihash.FromB58String(hex)
	if er != nil {
		fmt.Fprintf(os.Stderr, "ERROR in mhashFromHexString:  Got back err=%v\n", er)
	}

	return MHashFromMultihash(m), er == nil
}

// yawn. exact copy of sha1FromHexString.
func mhashFromHexBytes(hex []byte) (digestType, bool) {
	s := string(hex)
	return mhashFromHexString(s)
}

func MHashFromString(s string) Ref {
	s1 := sha256.New()
	s1.Write([]byte(s))
	d := mhashFromBinary(s1.Sum(nil))
	r := Ref{digest: d}
	return r //RefFromHash(s1)
}

func MHashFromMultihash(h multihash.Multihash) Mhash {
	var d Mhash
	if len(h) == 0 {
		return d
	}
	dec, er := multihash.Decode(h)
	if er != nil {
		fmt.Fprintf(os.Stderr, "!!!Decoded says format is %s and encoding %s\n", dec.Name, dec.Length)
		fmt.Fprintf(os.Stderr, "!!!There's an error as well: %v\n", er)
		print(er)
	}
	copy(d[:], h[2:])
	return d
}
