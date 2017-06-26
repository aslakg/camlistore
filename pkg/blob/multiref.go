package blob

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"os"

	"github.com/multiformats/go-multihash"
)

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
	m, er := multihash.FromB58String(hex)
	if er != nil {
		fmt.Fprintf(os.Stderr, "ERROR in mhashFromHexString: %v\n", er)
	}

	return mhashFromMultihash(m), er == nil
}

// yawn. exact copy of sha1FromHexString.
func mhashFromHexBytes(hex []byte) (digestType, bool) {
	s := string(hex)
	return mhashFromHexString(s)
}

func mhashFromString(s string) Ref {
	s1 := sha256.New()
	s1.Write([]byte(s))
	d := mhashFromBinary(s1.Sum(nil))
	r := Ref{digest: d}
	return r //RefFromHash(s1)
}

func mhashFromMultihash(h multihash.Multihash) Mhash {
	var d Mhash
	if len(h) == 0 {
		return d
	}
	dec, er := multihash.Decode(h)
	if er != nil {
		fmt.Fprintf(os.Stderr, "!!!Decoded says format is %s and encoding %s error: %v\n", dec.Name, dec.Length, er)
		print(er)
	}
	copy(d[:], h[2:])
	return d
}
